package responses

import "db_labs/entities"

type MeResponse struct {
	Name       string   `json:"name"`
	Surname    string   `json:"surname"`
	Patronymic string   `json:"patronymic"`
	Roles      []string `json:"roles"`
}

func MeFromEntity(usr *entities.User) *MeResponse {
	return &MeResponse{
		Name:       usr.Name(),
		Surname:    usr.Surname(),
		Patronymic: usr.Patronymic(),
		Roles:      usr.RoleNames(),
	}
}
