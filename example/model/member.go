package model

type Member struct {
	ID       int    `gorm:"column:id"`
	Username string `json:"username" form:"username" validate:"min=3,required" gorm:"column:username"`
	Email    string `json:"email" form:"email" validate:"email" gorm:"column:email"`
	Mobile   string `json:"mobile" form:"mobile" validate:"min=6" gorm:"column:mobile"`
	Nickname string `json:"nickname" form:"nickname" validate:"min=5,max=16" gorm:"column:nickname"`
	Avatar   string `json:"avatar" form:"avatar" validate:"url" gorm:"column:avatar"`
}

func (Member) TableName() string {
	return "user"
}
