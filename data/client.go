package data

/*
 * Struct defining identity of a client object.
 */
type Client struct {
	Name   string
	Role   string
	Status bool
}

// Returns Client map for the given struct.
func NewClient(name string, role string) *Client {
	return &Client{Name: name, Role: role, Status: true}
}
