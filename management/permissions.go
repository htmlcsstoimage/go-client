package management

// Permission identifies an API capability. Additional server permissions may be supplied as strings.
type Permission string

// API permissions available for grants and list filters.
const (
	PermissionAPIKeysCreateUpdate             Permission = "api_keys:create_update"
	PermissionAPIKeysDelete                   Permission = "api_keys:delete"
	PermissionAPIKeysRead                     Permission = "api_keys:read"
	PermissionImagesCreate                    Permission = "images:create"
	PermissionImagesDelete                    Permission = "images:delete"
	PermissionImagesRead                      Permission = "images:read"
	PermissionImagesStore                     Permission = "images:store"
	PermissionOGConfigsCreateUpdate           Permission = "og_configs:create_update"
	PermissionOGConfigsDelete                 Permission = "og_configs:delete"
	PermissionOGConfigsRead                   Permission = "og_configs:read"
	PermissionProxiesCreateUpdate             Permission = "proxies:create_update"
	PermissionProxiesDelete                   Permission = "proxies:delete"
	PermissionProxiesRead                     Permission = "proxies:read"
	PermissionStorageDestinationsCreateUpdate Permission = "storage_destinations:create_update"
	PermissionStorageDestinationsDelete       Permission = "storage_destinations:delete"
	PermissionStorageDestinationsRead         Permission = "storage_destinations:read"
	PermissionTemplatesCreateUpdate           Permission = "templates:create_update"
	PermissionTemplatesDelete                 Permission = "templates:delete"
	PermissionTemplatesRead                   Permission = "templates:read"
	PermissionUsageRead                       Permission = "usage:read"
)
