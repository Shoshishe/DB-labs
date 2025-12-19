package responses

type Role struct {
	Name string `json:"name"`
	Id   int8   `json:"id"`
}

func NewRole(Id int8, Name string) *Role {
	return &Role{Id: Id, Name: Name}
}
