package entity

type User struct {
	ID             int    `gorm:"primaryKey;column:id"`
	FirstName      string `gorm:"size:50;not null;column:first_name"`
	LastName       string `gorm:"size:50;unique;not null;column:last_name"`
	Username       string `gorm:"size:50;unique;not null;column:username"`
	Email          string `gorm:"size:100;unique;not null;column:email"`
	Password       string `gorm:"size:255;not null;column:password"`
	ProfilePicture string `gorm:"size:255;column:profile_picture"`
}

func (u *User) TableName() string {
	return "users"
}
