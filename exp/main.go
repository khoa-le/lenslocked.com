package main

import(
	"fmt"
	"bufio"
	"os"
	"strings"
	"github.com/jinzhu/gorm"
 	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	host = "localhost"
	port = "5432"
	user = "khoa"
	dbname = "lenslocked_dev"
)

type User struct{
	gorm.Model
	Name string
	Email string `gorm:"not null;unique_index"`
	Orders []Order

}

type Order struct{
	gorm.Model
	UserID uint
	Amount int
	Description string
}

func main(){
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
	host, port, user, dbname)
	db,err := gorm.Open("postgres", psqlInfo)
	if err != nil{
		panic(err)
	}

	defer db.Close()

	if err := db.DB().Ping(); err !=nil{
		panic(err)
	}	

	db.LogMode(true)
	db.AutoMigrate(&User{},&Order{})

	// name, email := getInfo()
	// u:=User{
	// 	Name: name,
	// 	Email: email,
	// }
	// db.Create(&u)
	// fmt.Println(u)

	// var u1 User
	// db.First(&u1)
	// fmt.Println(u1)

	// var users []User
	// db.Find(&users)
	// fmt.Println(len(users))
	// fmt.Println(users)

	var u User
	if err := db.Preload("Orders").First(&u).Error; err != nil{
		panic(err)
	}
	
	fmt.Println(u);
	//createOrder(db,u,1001,"Fake Description Order 1")

}

func getInfo() (name, email string){
	reader :=bufio.NewReader(os.Stdin)
	fmt.Println("What is your name?")
	name,_ = reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Println("What is your email?")
	email,_ = reader.ReadString('\n')
	email = strings.TrimSpace(email)
	return name,email;
}

func createOrder(db *gorm.DB,user User,amount int, desc string){
	err:=db.Create(&Order{
		UserID: user.ID,
		Amount: amount,
		Description: desc,
	}).Error
	if err != nil{
		panic(err)
	}
} 