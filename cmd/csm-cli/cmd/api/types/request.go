package types

// User - Struct for user
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"jwtToken,omitempty"`
}

// Cluster - Struct for cluster
type Cluster struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}

// Storage - Struct for storage array
type Storage struct {
	Endpoint    string `json:"management_endpoint,omitempty"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	UniqueID    string `json:"unique_id,omitempty"`
	StorageType string `json:"storage_array_type,omitempty"`
}
