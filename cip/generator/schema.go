package generator

// EventField represents a field in an event
type EventField struct {
	FieldName string `json:"name"`     // Field name
	Type      string `json:"type"`     // Field type
	Tag       string `json:"tag"`      // Tag name in Nostr event
	Required  bool   `json:"required"` // Whether this field is required
	Multiple  bool   `json:"multiple"` // Whether this field can have multiple values
}

// EventDefinition represents the definition of an event type
type EventDefinition struct {
	EventName   string       `json:"name"`        // Event name (e.g., "PostEvent")
	Operation   string       `json:"operation"`   // Operation type (e.g., "post")
	Kind        int          `json:"kind"`        // Event kind number
	Fields      []EventField `json:"fields"`      // Event fields
	Description string       `json:"description"` // Event description
}

// CIPDefinition represents a complete CIP implementation
type CIPDefinition struct {
	CIPName     string            `json:"name"`        // CIP name (e.g., "governance")
	Package     string            `json:"package"`     // Package name (e.g., "cip01")
	Description string            `json:"description"` // CIP description
	Events      []EventDefinition `json:"events"`      // List of events in this CIP
}
