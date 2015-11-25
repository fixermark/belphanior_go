package servant

// An individual handler for a Belphanior method
type Handler struct {
	Name    string      `json:"name"`
	Method  string      `json:"method"`
	Path    string      `json:"path"`
	Data    string      `json:"data,omitempty"`
	Handler interface{} `json:"-"`
}

// A full Belphanior role implementation.
type RoleImplementation struct {
	RoleUrl  string    `json:"role_url"`
	Handlers []Handler `json:"handlers"`
}
