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

//User represents the user model stored in our database
//This is used for user account
type User struct{
	gorm.Model
	Name string
	Email string `gorm:"not null;unique_index"`
	Password string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

//UserDB is used to interact with the users database
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

//UserService is a set of methods used to manipulate and work with the user model
type UserService interface{
	//Authenticate will verify the provided email and password are correct. If they
	//are correct, the user corresponding to that email will return. Otherwise You
	//will receive either: ErrNotFound, ErrInvalidPassword or other error if something 
	// goes wrong.
	Authenticate(email, password string) (*User, error)

	UserDB
}

func NewUserService(connectionInfo string) (UserService, error){
	ug, err := newUserGorm(connectionInfo)
	if err != nil{
		return nil, err
	}
	hmac := hash.NewHMAC(hmacSecretKey)
	uv := &userValidator{
		hmac: hmac,
		UserDB: ug,
	}
	return &userService{
		UserDB :uv,
	}, nil
}

var _ UserService = &userService{}

type userService struct{
	UserDB
}

//Authenticate can be used to authenticate a user with the email and password
func (us *userService) Authenticate(email,password string) (*User, error){
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

type userValidatorFunction func(*User) error

func runUserValidatorFunction(user *User, fns ...userValidatorFunction) error{
	for _,fn :=range(fns){
		if err := fn(user); err!=nil{
			return err
		}
	}
	return nil
}

var _ UserDB = &userValidator{}

type userValidator struct{
	UserDB
	hmac hash.HMAC
}

//ByRemember will hash the remember token and call the ByRemember
// on the subsequent on UserDB layer
func (uv *userValidator) ByRemember(token string) (*User, error){
	rememberHash := uv.hmac.Hash(token)
	return uv.UserDB.ByRemember(rememberHash)
}

//Create will create the provided user and backfill data
func (uv *userValidator) Create(user *User) error{
	if err:=runUserValidatorFunction(user, uv.bcryptPassword); err!=nil{
		return err;
	}

	if user.Remember == ""{
		token, err  := rand.RememberToken()
		if err != nil{
			return err
		}
		user.Remember = token
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)

	return uv.UserDB.Create(user)
}

//Update 
func (uv *userValidator) Update(user *User) error{
	if user.Remember != ""{
		user.RememberHash = uv.hmac.Hash(user.Remember)
	}
	return uv.UserDB.Update(user)
}

//Delete will delele the user with the provided ID
func (uv *userValidator) Delete(id uint) error{
	if id == 0{
		return ErrInvalidID
	}
	return uv.UserDB.Delete(id)
}

//bcryptPassword will hash a user's password with a predefined 
//pepper (userPepper) and bcrypt
func (uv *userValidator) bcryptPassword(user *User) error{
	if user.Password == ""{
		return nil
	}
	pwBytes := []byte(user.Password + userPwdPepper)
	hashedBytes,err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err!=nil{
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

var _ UserDB = &userGorm{}

func newUserGorm(connectionInfo string) (*userGorm, error){
	db,err := gorm.Open("postgres", connectionInfo)
	if err != nil{
		panic(err)
	}
	db.LogMode(true)
	return &userGorm{
		db:db,
	}, nil
}

type userGorm struct{
	db *gorm.DB
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
//This method expects the rememeber already to hashed.
func (ug *userGorm) ByRemember(rememberHash string) (*User, error){
	var user User
	err := first(ug.db.Where("remember_hash = ?", rememberHash),&user)
	if err != nil{
		return nil, err
	} 

	return &user, nil
}

//Create will create the provided user and backfill data
func (ug *userGorm) Create(user *User) error{
	return ug.db.Create(user).Error
}

//Update will update the provided user with all the data
//in provided user object 
func (ug *userGorm) Update(user *User) error{
	return ug.db.Save(user).Error
}

//Delete will delele the user with the provided ID
func (ug *userGorm) Delete(id uint) error{
	user :=User{Model:gorm.Model{ID:id}}
	return ug.db.Delete(&user).Error
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

func first(db *gorm.DB, dst interface{}) error{
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound{
		return ErrNotFound
	}
	return err
}
