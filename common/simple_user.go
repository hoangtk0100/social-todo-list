package common

type SimpleUser struct {
	SQLModel
	FirstName string `json:"first_name" gorm:"column:first_name;"`
	LastName  string `json:"last_name" gorm:"column:last_name;"`
	Status    int    `json:"status" gorm:"column:status;"`
}

func (SimpleUser) TableName() string {
	return "users"
}

func (s *SimpleUser) Mask() {
	s.SQLModel.Mask(DBTypeUser)
}
