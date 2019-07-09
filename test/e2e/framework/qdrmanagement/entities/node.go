package entities

// Node represents a Dispatch Router router.node entity
type Node struct {
	Id              string   `json:"id"`
	ProtocolVersion int      `json:"protocolVersion"`
	Instance        int      `json:"instance"`
	LinkState       []string `json:"linkState"`
	NextHop         string   `json:"nextHop"`
	ValidOrigins    []string `json:"validOrigins"`
	Address         string   `json:"address"`
	RouterLink      int      `json:"routerLink"`
	Cost            int      `json:"cost"`
	LastTopoChange  int      `json:"lastTopoChange"`
	Index           int      `json:"index"`
}

// Implementation of the Entity interface
func (Node) GetEntityId() string {
	return "router.node"
}
