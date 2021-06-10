package domain

type VmLogin struct {
	VmName     string `bson:"vm_name" json:"virtual_machine_name"`
	VmUserName string `bson:"vm_username" json:"username,omitempty"`
	VmPassword string `bson:"vm_password" json:"password,omitempty"`
	IpAdd      string `bson:"vm_ip" json:"public_ip_address"`
}
type ResourceGroup struct {
	Name     string    `bson:"resourcegroup_name" json:"resourcegroup_name"`
	Region   string    `bson:"region" json:"region"`
	LoginDet []VmLogin `bson:"virtual_machine" json:"vm_credentials"`
}

type Organization struct {
	OrgName       string          `bson:"org_name" json:"organization_name"`
	ResourceGroup []ResourceGroup `bson:"resourcegroup" json:"resourcegroups"`
}

type VmAll struct {
	OrgName    string `bson:"org_name" json:"orgnaization_name"`
	ResGrpName string `bson:"resourcegroup_name" json:"resourcegroup_name"`
	Region     string `bson:"region" json:"region"`
	VmName     string `bson:"vm_name"`
	VmUserName string `bson:"vm_username"`
	VmPassword string `bson:"vm_password"`
	IpAdd      string `bson:"vm_ip"`
}

type Location struct {
	// ResourcrGroupName string `bson:"resourcegroup"`
	Region string `bson:"region"`
}
