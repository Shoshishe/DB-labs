package ioc

import "db_labs/controllers"

var useAuthController = provider(
	func() *controllers.AuthController {
		return controllers.NewAuthController(UseHttpMux(), useAuthService())
	},
)
