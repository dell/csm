package api

const (
	// UserLoginURI - User login URI
	UserLoginURI = "/api/v1/users/login"

	// AddCLusterURI - Add CLuster URI
	AddCLusterURI = "/api/v1/clusters"

	// GetAllClustersURI - Get All Clusters URI
	GetAllClustersURI = "/api/v1/clusters"

	// GetClusterByNameURI - Get Cluster by name URI
	GetClusterByNameURI = GetAllClustersURI + "?cluster_name=%s"

	// PatchClusterURI - Patch Cluster URI
	//@TODO - check if /api/v1/clusters/name:%s can be used
	PatchClusterURI = "/api/v1/clusters/%d"

	// DeleteClusterURI - Delete Cluster URI
	DeleteClusterURI = "/api/v1/clusters/%d"

	// AddStorageURI - Add Storage array URI
	AddStorageURI = "/api/v1/storage-arrays"

	// GetAllStorageURI - Get All Storage arrays URI
	GetAllStorageURI = "/api/v1/storage-arrays"

	// GetStorageByParamURI - Get Storage array by custom parameter URI
	GetStorageByParamURI = GetAllStorageURI + "?%s=%s"

	// DeleteStorageURI - Delete Storage array URI
	DeleteStorageURI = "/api/v1//storage-arrays/%d"
)
