package ioc

import "db_labs/services"

var useAuthService = provider(func() *services.AuthService { return services.NewAuthService(useAuthRepo()) })

var useUniversityService = provider(func() *services.UniversityService {
	return services.NewUniversityService(UseUniversitiesRepo())
})

var useUserService = provider(func() *services.UsersService {
	return services.NewUsersService(useUsersRepo())
})

var useGroupsService = provider(func() *services.GroupsService {
	return services.NewGroupsService(useGroupsRepo())
})
