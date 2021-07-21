package api

const (
	// User login URI
	UserLoginURI = "/api/v1/users/login"

	// Add CLuster URI
	AddCLusterURI = "/api/v1/clusters"

	// Get All Clusters URI
	GetAllClustersURI = "/api/v1/clusters"

	// Get Cluster by name URI
	GetClusterByNameURI = GetAllClustersURI + "?cluster_name=%s"

	// Patch Cluster URI
	//@TODO - check if /api/v1/clusters/name:%s can be used
	PatchClusterURI = "/api/v1/clusters/%d"

	// Delete Cluster URI
	DeleteClusterURI = "/api/v1/clusters/%d"

	// Add Storage array URI
	AddStorageURI = "/api/v1/storage-arrays"

	// Get All Storage arrays URI
	GetAllStorageURI = "/api/v1/storage-arrays"

	// Get Storage array by custom parameter URI
	GetStorageByParamURI = GetAllStorageURI + "?%s=%s"

	// Delete Storage array URI
	DeleteStorageURI = "/api/v1//storage-arrays/%d"
)
