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
	db.LogMode(true)
	return &UserService{
		db:db,
	}, nil
}

type UserService struct{
	db *gorm.DB
}

var (
	ErrNotFound = errors.New("models: resource not found")

	ErrInvalidID = errors.New("models: ID provided was invalid")
)

//ByID will look up a user with the provided  ID.
func (us *UserService) ByID(id uint) (*User,error){
	var user User
	db := us.db.Where("id=?",id)
	err := first(db,&user)
	return &user, err
}

//Create will create the provided user and backfill data
func (us *UserService) Create(user *User) error{
	return us.db.Create(user).Error
}

//Update user
func (us *UserService) Update(user *User) error{
	return us.db.Save(user).Error
}

//ByEmail will looks up a user with a given email address and return that user
func (us *UserService) ByEmail(email string) (*User,error){
	var user User
	db := us.db.Where("email=?",email)
	err := first(db,&user)
	return &user, err
}

func (us *UserService) Delete(id uint) error{
	if id == 0{
		return ErrInvalidID
	}
	user :=User{Model:gorm.Model{ID:id}}
	return us.db.Delete(&user).Error
}

func first(db *gorm.DB, dst interface{}) error{
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound{
		return ErrNotFound
	}
	return err
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