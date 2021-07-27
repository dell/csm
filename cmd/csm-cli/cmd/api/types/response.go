package types

// JWTToken - captures JWT Token
var JWTToken string

// ClusterResponse - Struct to capture cluster response
type ClusterResponse struct {
	ClusterID   uint   `json:"cluster_id"`
	ClusterName string `json:"cluster_name"`
	Nodes       string `json:"nodes"`
}

// StorageResponse - Struct to capture storage array response
type StorageResponse struct {
	ID            uint   `json:"id"`
	StorageTypeID uint   `json:"storage_array_type_id"`
	UniqueID      string `json:"unique_id"`
	Username      string `json:"username"`
	Endpoint      string `json:"management_endpoint"`
}
