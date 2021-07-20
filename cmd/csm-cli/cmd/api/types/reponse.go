package types

// Struct to capture user login response
type UserLoginResponse struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

// Struct to capture cluster response
type ClusterResponse struct {
	ClusterId   uint   `json:"cluster_id"`
	ClusterName string `json:"cluster_name"`
	Nodes       string `json:"nodes"`
}

// Struct to capture storage array response
type StorageResponse struct {
	Id            uint   `json:"id"`
	StoragtTypeId uint   `json:"storage_array_type_id"`
	UniqueId      string `json:"unique_id"`
	Username      string `json:"username"`
	Endpoint      string `json:"management_endpoint"`
}
