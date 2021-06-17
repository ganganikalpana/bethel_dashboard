package domain

import "go.mongodb.org/mongo-driver/bson/primitive"



type Metrics struct {
	MetricId          primitive.ObjectID     `bson:"_id,omitempty"`
	VMID              interface{}            `bson:"virtualmachine,omitempty"`
	NetworkIn         map[string]interface{} `bson:"network_in_total"`
	NetworkOut        map[string]interface{} `bson:"network_out_total"`
	PercentCpu        map[string]interface{} `bson:"percantage_cpu"`
	AvailableMemBytes map[string]interface{} `bson:"available_memory_bytes"`
}
type VirtualMachine struct {
	VMID       primitive.ObjectID `bson:"_id,omitempty"`
	ResGrpId   interface{}        `bson:"resourcegroup,omitempty"`
	VmName     string             `bson:"vm_name" json:"virtual_machine_name"`
	VmUserName string             `bson:"vm_username" json:"username,omitempty"`
	VmPassword string             `bson:"vm_password" json:"password,omitempty"`
	IpAdd      string             `bson:"vm_ip" json:"public_ip_address"`
}
type ResourceGroup struct {
	ResGrpId primitive.ObjectID `bson:"_id,omitempty"`
	OrgId    interface{}        `bson:"organization,omitempty"`
	Name     string             `bson:"resourcegroup_name" json:"resourcegroup_name"`
	Region   string             `bson:"region" json:"region"`
}

type Organization struct {
	OrgId   primitive.ObjectID `bson:"_id,omitempty"`
	OrgName string             `bson:"org_name" json:"organization_name"`
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
	Region string `bson:"region"`
}
