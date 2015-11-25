package servant

// A full protocol implementation for a servant
type Protocol struct {
	Roles []RoleImplementation `json:"roles"`
}
