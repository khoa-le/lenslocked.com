package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"golang.org/x/crypto/bcrypt"

	"lenslocked.com/hash"
	"lenslocked.com/rand"

)

var (
	ErrNotFound = errors.New("models: resource not found")
	ErrInvalidID = errors.New("models: ID provided was invalid")
	ErrInvalidPassword = errors.New("models: invalid password prodvided")
)
const userPwdPepper = "secret-random-string"
const hmacSecretKey = "secret-hmac-key"

type UserDB interface{
	//Query for single user
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	//Methods for altering user
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	//Use to close db
	Close() error

	//Migration helper
	AutoMigrate() error
	DestructiveReset() error
}

func NewUserService(connectionInfo string) (*UserService, error){
	ug, err := newUserGorm(connectionInfo)
	if err != nil{
		return nil, err
	}
	return &UserService{
		UserDB : &userValidator{
			UserDB: ug,
		},
	}, nil
}

type UserService struct{
	UserDB
}

type userValidator struct{
	UserDB
}

func newUserGorm(connectionInfo string) (*userGorm, error){
	db,err := gorm.Open("postgres", connectionInfo)
	if err != nil{
		panic(err)
	}
	db.LogMode(true)
	hmac := hash.NewHMAC(hmacSecretKey)
	return &userGorm{
		db:db,
		hmac: hmac,
	}, nil
}
var _ UserDB = &userGorm{}

type userGorm struct{
	db *gorm.DB
	hmac hash.HMAC
}

//ByID will look up a user with the provided  ID.
func (ug *userGorm) ByID(id uint) (*User,error){
	var user User
	db := ug.db.Where("id=?",id)
	err := first(db,&user)
	return &user, err
}

//ByEmail will looks up a user with a given email address and return that user
func (ug *userGorm) ByEmail(email string) (*User,error){
	var user User
	db := ug.db.Where("email=?",email)
	err := first(db,&user)
	return &user, err
}

//ByRemember will looks up a user with a given token string and return that user
//This method will handle hashing token for us
func (ug *userGorm) ByRemember(token string) (*User, error){
	var user User
	rememberHash := ug.hmac.Hash(token)
	err := first(ug.db.Where("remember_hash = ?", rememberHash),&user)
	if err != nil{
		return nil, err
	}

	return &user, nil;
}

//Create will create the provided user and backfill data
func (ug *userGorm) Create(user *User) error{
	pwBytes := []byte(user.Password + userPwdPepper)
	hashedBytes,err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err!=nil{
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	
	if user.Remember == ""{
		token, err  := rand.RememberToken()
		if err != nil{
			return err
		}
		user.Remember = token
	}
	user.RememberHash = ug.hmac.Hash(user.Remember)

	return ug.db.Create(user).Error
}

//Update user
func (ug *userGorm) Update(user *User) error{
	if user.Remember != ""{
		user.RememberHash = ug.hmac.Hash(user.Remember)
	}
	return ug.db.Save(user).Error
}

func (ug *userGorm) Delete(id uint) error{
	if id == 0{
		return ErrInvalidID
	}
	user :=User{Model:gorm.Model{ID:id}}
	return ug.db.Delete(&user).Error
}


//Authenticate can be used to authenticate a user with the email and password
func (us *UserService) Authenticate(email,password string) (*User, error){
	foundUser, err := us.ByEmail(email)
	if err != nil{
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash),[]byte(password + userPwdPepper))
	if err != nil{
		switch err {
			case bcrypt.ErrMismatchedHashAndPassword: 
				return nil, ErrInvalidPassword
			default:
				return nil, err
		}
	} 
	return foundUser, nil

}

func first(db *gorm.DB, dst interface{}) error{
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound{
		return ErrNotFound
	}
	return err
}

//Closes the UserService database connection 
func (ug *userGorm) Close() error{
	return ug.db.Close()
}

//DestructiveReset drops the user table and rebuild it
func (ug *userGorm) DestructiveReset() error{

	if err:=ug.db.DropTableIfExists(&User{}).Error; err !=nil{
		return err
	}

	return ug.AutoMigrate()
}

func (ug *userGorm) AutoMigrate() error{
	if err := ug.db.AutoMigrate(&User{}).Error;err!=nil{
		return err
	}
	return nil
}

type User struct{
	gorm.Model
	Name string
	Email string `gorm:"not null;unique_index"`
	Password string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}