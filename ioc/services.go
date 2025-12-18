package ioc

import "db_labs/services"

var useAuthService = provider(func() *services.AuthService { return services.NewAuthService(useAuthRepo()) })
