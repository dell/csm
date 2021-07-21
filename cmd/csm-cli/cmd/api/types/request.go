package types

// Struct for user
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"jwtToken,omitempty"`
}

// Struct for cluster
type Cluster struct {
	Name string `json:"name,omitempty"`
	Id   string `json:"id,omitempty"`
}

// Struct for storage array
type Storage struct {
	Endpoint    string `json:"management_endpoint,omitempty"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	UniqueId    string `json:"unique_id,omitempty"`
	StorageType string `json:"storage_array_type,omitempty"`
}
