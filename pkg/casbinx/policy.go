package casbinx

func Policies(service *CasbinService) {
	// Admin can manage roles
	service.AddPermission("admin", "/v1/role/create", "POST")
	service.AddPermission("admin", "/v1/role/update/:id", "PUT")
	service.AddPermission("admin", "/v1/role/:id", "DELETE")

}
