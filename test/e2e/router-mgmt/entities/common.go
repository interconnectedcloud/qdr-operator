package entities

// EntityCommon represents the common attributes
// of any Dispatch Router Entity
type EntityCommon struct {
	Name     string `json:"name"`
	Identity string `json:"identity"`
	Type     string `json:"type"`
}

// Entity interface must be implemented by all entities
type Entity interface {
	GetEntityId() string
}
