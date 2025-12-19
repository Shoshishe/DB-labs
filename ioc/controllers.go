package ioc

import (
	"db_labs/controllers"
	"net/http"
)

var UseAuthController = provider(
	func() *controllers.AuthController {
		return controllers.NewAuthController(UseHttpMux(), useAuthService())
	},
)

var UseUniversitiesController = provider(
	func() *controllers.UniversityController {
		return controllers.NewUniversityController(UseHttpMux(), useUniversityService())
	},
)

var UseUsersController = provider(
	func() *controllers.UserController {
		return controllers.NewUserController(useUserService())
	},
)

var UseGroupsController = provider(
	func() *controllers.GroupsController {
		return controllers.NewGroupsController(useGroupsService())
	},
)

type Controller interface {
	RegisterRoutes(mux *http.ServeMux)
}

func SetupControllers(controllers []Controller) {
	for _, controller := range controllers {
		controller.RegisterRoutes(UseHttpMux())
	}
}
