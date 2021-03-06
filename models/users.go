package models

import (
	"strings"
	"regexp"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"golang.org/x/crypto/bcrypt"

	"lenslocked.com/hash"
	"lenslocked.com/rand"

)

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
}

//UserService is a set of methods used to manipulate and work with the user model
type UserService interface{
	//Authenticate will verify the provided email and password are correct. If they
	//are correct, the user corresponding to that email will return. Otherwise You
	//will receive either: ErrNotFound, ErrPasswordIncorrect or other error if something 
	// goes wrong.
	Authenticate(email, password string) (*User, error)

	UserDB
}

func NewUserService(db *gorm.DB, pepper string, hmacKey string) UserService{
	ug := &userGorm{db}
	hmac := hash.NewHMAC(hmacKey)

	uv := newUserValidator(ug, hmac, pepper)

	return &userService{
		UserDB :uv,
		pepper: pepper,
	}
}

var _ UserService = &userService{}

type userService struct{
	UserDB
	pepper string
}

//Authenticate can be used to authenticate a user with the email and password
func (us *userService) Authenticate(email,password string) (*User, error){
	foundUser, err := us.ByEmail(email)
	if err != nil{
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash),[]byte(password + us.pepper))
	if err != nil{
		switch err {
			case bcrypt.ErrMismatchedHashAndPassword: 
				return nil, ErrPasswordIncorrect
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

func newUserValidator(udb UserDB, hmac hash.HMAC, pepper string) *userValidator{
	return &userValidator{
		UserDB: udb,
		hmac: hmac,
		emailRegex: regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"),
		pepper: pepper,
	}
}

type userValidator struct{
	UserDB
	hmac hash.HMAC
	emailRegex *regexp.Regexp
	pepper string
}

//ByEmail will normalize the email address before calling
//ByEmail on the DB layer
func (uv *userValidator) ByEmail(email string) (*User, error){
	user := User{
		Email: email,
	}
	if err := runUserValidatorFunction(&user, uv.normalizeEmail); err!=nil{
		return nil, err
	}
	return uv.UserDB.ByEmail(user.Email)
}

//ByRemember will hash the remember token and call the ByRemember
// on the subsequent on UserDB layer
func (uv *userValidator) ByRemember(token string) (*User, error){
	var user = User{
		Remember: token,
	}
	if err := runUserValidatorFunction(&user, uv.hmacRemember); err!=nil{
		return nil, err
	}
	
	return uv.UserDB.ByRemember(user.RememberHash)
}

//Create will create the provided user and backfill data
func (uv *userValidator) Create(user *User) error{
	if err:=runUserValidatorFunction(user, 
		uv.passwordRequired,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.setRememberIfUnset, 
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvailable); err!=nil{
		return err
	}

	return uv.UserDB.Create(user)
}

//Update 
func (uv *userValidator) Update(user *User) error{
	if err:=runUserValidatorFunction(user, 
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvailable); err!=nil{
		return err
	}

	return uv.UserDB.Update(user)
}

//Delete will delele the user with the provided ID
func (uv *userValidator) Delete(id uint) error{
	var user User
	user.ID = id
	err := runUserValidatorFunction(&user, uv.idGreaterThan(0))
	if err != nil{
		return err
	}
	return uv.UserDB.Delete(id)
}

//bcryptPassword will hash a user's password with a predefined 
//pepper (userPepper) and bcrypt
func (uv *userValidator) bcryptPassword(user *User) error{
	if user.Password == ""{
		return nil
	}
	pwBytes := []byte(user.Password + uv.pepper)
	hashedBytes,err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err!=nil{
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

func (uv *userValidator) hmacRemember(user *User) error{
	if user.Remember == ""{
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)

	return nil
}

func (uv *userValidator) setRememberIfUnset(user *User) error{
	if user.Remember != ""{
		return nil
	}
	token, err  := rand.RememberToken()
	if err != nil{
		return err
	}
	user.Remember = token
	return nil
}
func (uv *userValidator) rememberMinBytes(user *User) error{
	if user.Remember == ""{
		return nil
	}
	n,err := rand.NBytes(user.Remember)
	if err != nil{
		return err
	}
	if n < 32{
		return ErrRememberTooShort
	}

	return nil
}

func (uv *userValidator) rememberHashRequired(user *User) error{
	if user.RememberHash == ""{
		return ErrRememberHashRequired
	}
	return nil
}

func (uv *userValidator) idGreaterThan(n uint) userValidatorFunction{
	return userValidatorFunction(func(user *User) error{
		if user.ID<= n{
			return ErrIDInvalid
		}
		return nil
	})
}

func (uv *userValidator) normalizeEmail(user *User) error{
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

func (uv *userValidator) requireEmail(user *User) error{
	if user.Email == "" {
		return ErrEmailRequired
	}
	return nil
}

func (uv *userValidator) emailFormat(user *User) error{
	if user.Email != "" && !uv.emailRegex.MatchString(user.Email){
		return ErrEmailInvalid
	}
	return nil
}

func (uv *userValidator) emailIsAvailable(user *User) error{
	existing, err := uv.ByEmail(user.Email)
	if err == ErrNotFound{
		return nil
	}
	if err !=nil{
		return err
	}

	if user.ID != existing.ID{
		return ErrEmailTaken
	}

	return nil
}

func (uv *userValidator) passwordMinLength(user *User) error{
	if user.Password == ""{
		return nil
	}

	if len(user.Password)<8{
		return ErrPasswordTooShort
	}
	return nil  
}

func (uv *userValidator) passwordRequired(user *User) error{
	if user.Password == ""{
		return ErrPasswordRequired
	}
	return nil
}

func (uv *userValidator) passwordHashRequired(user *User) error{
	if user.PasswordHash == ""{
		return ErrPasswordHashRequired
	}
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

//Delete will delete the user with the provided ID
func (ug *userGorm) Delete(id uint) error{
	user :=User{Model:gorm.Model{ID:id}}
	return ug.db.Delete(&user).Error
}

func first(db *gorm.DB, dst interface{}) error{
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound{
		return ErrNotFound
	}
	return err
}
