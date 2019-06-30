package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func NewUserService(connectionInfo string) (*UserService, error){
	db,err := gorm.Open("postgres", connectionInfo)
	if err != nil{
		panic(err)
	}

	return &UserService{
		db:db,
	}, nil
}

type UserService struct{
	db *gorm.DB
}

var (
	ErrNotFound = errors.New("models: resource not found")
)

//ByID will look up a user with the provided  ID.
func (us *UserService) ByID(id uint) (*User,error){
	var user User
	err := us.db.Where("id=?",id).First(&user).Error
	switch err {
	case nil:
		return &user, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

//Closes the UserService database connection 
func (us *UserService) Close() error{
	return us.db.Close()
}

//DestructiveReset drops the user table and rebuild it
func (us *UserService) DestructiveReset(){
	us.db.DropTableIfExists(&User{})
	us.db.AutoMigrate(&User{})
}



type User struct{
	gorm.Model
	Name string
	Email string `gorm:"not null;unique_index"`
}