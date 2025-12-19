package ioc

import (
	"db_labs/repository"
)

var useAuthRepo = provider(
	func() *repository.AuthRepistory {
		return repository.NewAuthRepository(*useTokenStore(), *useUsersRepo())
	},
)

