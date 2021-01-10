package user

type User struct {
	Uid  int    `gorm:"not null;PRIMARY_KEY;comment:'用户uid'" json:"uid"` // 用户
	Name string `gorm:"type:json;comment:'名字'" json:"name"`              // 名字
}
