package model

// Role permission rules
type RoleCasbin struct {
	Keyword string `json:"keyword"` // Role keyword
	Path    string `json:"path"`    // Access path
	Method  string `json:"method"`  // Request method
}
