package dto

import "github.com/niluwats/bethel_dashboard/domain"

type NewNodeRequest struct {
	OrgName        string `json:"organization"`
	ResGroup       string `json:"resource_group"`
	NodeName       string `json:"node_name"`
	NodeNIC        string `json:"node_nic"`
	NodeIp         string `json:"nodeip"`
	NodeDeployment string `json:"node_deployment"`
	Region         string `json:"region"`
}

type NewNodesResponse struct {
	ResGroup []domain.ResourceGroup `json:"resourcegroups"`
}

func ToDto(org *domain.Organization) *NewNodesResponse {
	nodeResponse := NewNodesResponse{
		ResGroup: org.ResourceGroup,
	}
	return &nodeResponse
}
