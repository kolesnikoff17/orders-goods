package entity

// Good -.
type Good struct {
	ID         string                 `json:"-" bson:"-"`
	Name       string                 `json:"name" bson:"name"`
	Category   string                 `json:"category" bson:"category"`
	Price      string                 `json:"price" bson:"price"`
	Additional map[string]interface{} `json:"additional,omitempty" bson:"additional,omitempty"`
}
