package dto

type NewNodeRequest struct {
	OrgName        string `json:"organization"`
	ResGroup       string `json:"resource_group"`
	NodeName       string `json:"node_name"`
	NodeNIC        string `json:"node_nic"`
	NodeIp         string `json:"nodeip"`
	NodeDeployment string `json:"node_deployment"`
	Region         string `json:"region"`
}

// type NewNodesRequest struct {
// 	ResGroup string `json:"resource_group"`
// 	Region   string `json:"region"`
// 	NodesArr []NewNode
// }

type NewNodesResponse struct {
	ResGroup       string `json:"resource_group"`
}
