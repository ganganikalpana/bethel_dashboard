package azure_gosdk

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/niluwats/bethel_dashboard/domain"
	"github.com/niluwats/bethel_dashboard/dto"
)

type Value struct {
	Vmname, Ip, NIC, deploymentName string
}

var (
	m                     = map[string]Value{}
	VmDetails             domain.VmAll
	organizationName      = ""
	resourceGroupName     = ""
	resourceGroupLocation = ""
	templateFile          = "azure_gosdk/quickstarts/deploy-vm/vm-quickstart-template.json"
)

func AssignMaps(dtoMap dto.NewNodeRequest) *domain.VmAll {
	m = map[string]Value{
		dtoMap.NodeName: {
			Vmname:         dtoMap.NodeName,
			Ip:             dtoMap.NodeIp,
			NIC:            dtoMap.NodeNIC,
			deploymentName: dtoMap.NodeDeployment,
		},
	}
	resourceGroupLocation = dtoMap.Region
	resourceGroupName = dtoMap.ResGroup
	organizationName = dtoMap.OrgName

	main()
	return &VmDetails
}

type mapCounter struct {
	mc map[string]Value
	sync.RWMutex
}

// Information loaded from the authorization file to identify the client
type clientInfo struct {
	SubscriptionID string
	VMPassword     string
}

var (
	ctx        = context.Background()
	clientData clientInfo
	authorizer autorest.Authorizer
)

// Authenticate with the Azure services using file-based authentication
func init() {
	var err error
	fmt.Println(os.Getenv("AZURE_AUTH_LOCATION"))
	authorizer, err = auth.NewAuthorizerFromFile(azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		log.Fatalf("Failed to get OAuth config: %v", err)
	}

	authInfo, err := readJSON(os.Getenv("AZURE_AUTH_LOCATION"))
	if err != nil {
		log.Fatalf("Failed to read JSON: %+v", err)
	}
	clientData.SubscriptionID = (*authInfo)["subscriptionId"].(string)
	clientData.VMPassword = (*authInfo)["clientSecret"].(string)
}

func main() {
	mc := mapCounter{
		mc: make(map[string]Value),
	}
	var wg sync.WaitGroup

	wg.Add(len(m))
	for _, v := range m {

		go func(v Value) {
			NewVm(&v, &mc)
			wg.Done()
		}(v)

	}
	wg.Wait()
	time.Sleep(3 * time.Second)

}
func NewVm(v *Value, mc *mapCounter) {
	fmt.Println("Starting")
	group, err := createGroup()
	if err != nil {
		log.Fatalf("failed to create group: %v", err)
	}
	log.Printf("Created group: %v", *group.Name)

	log.Printf("Starting deployment: %s", v.deploymentName)
	result, err := createDeployment(v, mc)
	if err != nil {
		log.Fatalf("Failed to deploy: %v", err)
	}
	if result.Name != nil {
		log.Printf("Completed deployment %v: %v", v.deploymentName, *result.Properties.ProvisioningState)
	} else {
		log.Printf("Completed deployment %v (no data returned to SDK)", v.deploymentName)
	}
	getLogin(v, mc)
	time.Sleep(2 * time.Second)
}

// Create a resource group for the deployment.
func createGroup() (group resources.Group, err error) {
	groupsClient := resources.NewGroupsClient(clientData.SubscriptionID)
	groupsClient.Authorizer = authorizer

	return groupsClient.CreateOrUpdate(
		ctx,
		resourceGroupName,
		resources.Group{
			Location: to.StringPtr(resourceGroupLocation)})
}

// Create the deployment
func createDeployment(v *Value, mc *mapCounter) (deployment resources.DeploymentExtended, err error) {

	template, err := readJSON(templateFile)
	if err != nil {
		return
	}
	param := Params(v, mc)

	deploymentsClient := resources.NewDeploymentsClient(clientData.SubscriptionID)
	deploymentsClient.Authorizer = authorizer
	c := &param

	deploymentFuture, err := deploymentsClient.CreateOrUpdate(
		ctx,
		resourceGroupName,
		v.Vmname,
		resources.Deployment{
			Properties: &resources.DeploymentProperties{
				Template:   template,
				Parameters: c,
				Mode:       resources.Incremental,
			},
		},
	)
	if err != nil {
		return
	}
	err = deploymentFuture.WaitForCompletionRef(ctx, deploymentsClient.BaseClient.Client)
	if err != nil {
		return
	}
	return deploymentFuture.Result(deploymentsClient)
}

// Get login information by querying the deployed public IP resource.
func getLogin(v *Value, mc *mapCounter) {
	param := Params(v, mc)
	addressClient := network.NewPublicIPAddressesClient(clientData.SubscriptionID)
	addressClient.Authorizer = authorizer
	ipName := param["publicIPAddresses_QuickstartVM_ip_name"].(map[string]interface{})

	ipAddress, err := addressClient.Get(ctx, resourceGroupName, ipName["value"].(string), "")
	if err != nil {
		log.Fatalf("Unable to get IP information. Try using `az network public-ip list -g %s", resourceGroupName)
	}

	vmUser := param["vm_user"].(map[string]interface{})
	vmName := param["virtualMachines_QuickstartVM_name"].(map[string]interface{})

	log.Printf("Log in with ssh: %s@%s, password: %s",
		vmUser["value"].(string),
		*ipAddress.PublicIPAddressPropertiesFormat.IPAddress,
		clientData.VMPassword)
	VmDetails = domain.VmAll{
		OrgName:    organizationName,
		ResGrpName: resourceGroupName,
		Region:     resourceGroupLocation,
		VmName:     vmName["value"].(string),
		VmPassword: clientData.VMPassword,
		IpAdd:      *ipAddress.IPAddress,
		VmUserName: vmUser["value"].(string),
	}

}

func readJSON(path string) (*map[string]interface{}, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}
	contents := make(map[string]interface{})
	_ = json.Unmarshal(data, &contents)
	return &contents, nil
}

func Params(v *Value, mc *mapCounter) map[string]interface{} {

	mc.Lock()
	var Param = map[string]interface{}{
		"virtualNetworks_GoQSVM_vnet_name":            map[string]interface{}{"value": "QuickstartVnet"},
		"virtualMachines_QuickstartVM_name":           map[string]interface{}{"value": v.Vmname},
		"networkInterfaces_quickstartvm_name":         map[string]interface{}{"value": v.NIC},
		"publicIPAddresses_QuickstartVM_ip_name":      map[string]interface{}{"value": v.Ip},
		"networkSecurityGroups_QuickstartVM_nsg_name": map[string]interface{}{"value": "QuickstartNSG"},
		"subnets_default_name":                        map[string]interface{}{"value": "QuickstartSubnet"},
		"securityRules_default_allow_ssh_name":        map[string]interface{}{"value": "qsuser"},
		"osDisk_name":                                 map[string]interface{}{"value": "_OsDisk_1_2e3ae1ad37414eaca81b432401fcdd75"},
		"vm_user":                                     map[string]interface{}{"value": "quickstart"},
		"vm_password":                                 map[string]interface{}{"value": "_"},
	}
	Param["vm_password"] = map[string]string{
		"value": clientData.VMPassword,
	}
	mc.Unlock()
	return Param
}
