package models

import "encoding/json"

// Bundle is documented here http://hl7.org/fhir/StructureDefinition/Bundle
// A container for a collection of resources.
type Bundle struct {
	ResourceType string `bson:"resourceType" json:"resourceType"`

	Id        string        `bson:"id,omitempty" json:"id,omitempty"`
	Type      string        `bson:"type" json:"type"`
	Timestamp string        `bson:"timestamp,omitempty" json:"timestamp,omitempty"`
	Total     int           `bson:"total,omitempty" json:"total,omitempty"`
	Entry     []BundleEntry `bson:"entry,omitempty" json:"entry,omitempty"`
}

// An entry in a bundle resource - will either contain a resource or information about a resource (transactions and history only).
type BundleEntry struct {
	FullUrl  string          `bson:"fullUrl,omitempty" json:"fullUrl,omitempty"`
	Resource json.RawMessage `bson:"resource,omitempty" json:"resource,omitempty"`
}
