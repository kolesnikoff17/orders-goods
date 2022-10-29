package entity

// Good -.
type Good struct {
  ID   string `json:"-"`
  Data map[string]interface{}
}
