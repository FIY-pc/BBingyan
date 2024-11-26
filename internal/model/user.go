package model

import (
	"github.com/FIY-pc/BBingyan/internal/config"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Email      string    `json:"email" gorm:"unique"`
	Nickname   string    `json:"nickname"`
	Password   string    `json:"password"`
	Permission int       `json:"permission" gorm:"default:1"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt"`

	Intro   string    `json:"intro"`
	Avatar  string    `json:"avatar"`
	Article []Article ` gorm:"foreignKey:UserID;-"` // forbid preload

	Nodes []Node `gorm:"many2many:user_nodes;" json:"nodes"` // 关联节点表
}

func InitUser(DB *gorm.DB) {
	if err := DB.AutoMigrate(&User{}); err != nil {
		panic(err)
	}
}

// InitSuperAdmin 初始化一个超级管理员
func InitSuperAdmin() {
	var err error
	_, err = GetUserByEmail(config.Config.User.Admin.Email)
	if err != nil {
		admin := &User{
			Email:      config.Config.User.Admin.Email,
			Password:   config.Config.User.Admin.Password,
			Nickname:   "",
			Permission: 10,
		}
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		admin.Password = string(hashPassword)
		err = CreateUser(admin)
		if err != nil {
			panic(err)
		}
	}
}

func CreateUser(user *User) error {
	if err := postgresDb.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func GetUserByEmail(email string) (*User, error) {
	var user User
	if err := postgresDb.Where("email", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByID(id uint) (*User, error) {
	var user User
	if err := postgresDb.Where("id", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(user *User) error {
	if err := postgresDb.Updates(user).Error; err != nil {
		return err
	}
	return nil
}

func DeleteUserByID(id uint) error {
	if err := postgresDb.Where("id", id).Delete(&User{}).Error; err != nil {
		return err
	}
	return nil
}

func DeleteUserByEmail(email string) error {
	if err := postgresDb.Where("email", email).Delete(&User{}).Error; err != nil {
		return err
	}
	return nil
}

func GetUsersByIDs(ids []uint) ([]User, error) {
	var users []User
	if err := postgresDb.Where("id IN ?", ids).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
