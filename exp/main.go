package main

import(
	"fmt"
	 _ "github.com/jinzhu/gorm/dialects/postgres"
	 
	 "lenslocked.com/models"
)

const (
	host = "localhost"
	port = "5432"
	user = "khoa"
	dbname = "lenslocked_dev"
)

func main(){
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
	host, port, user, dbname)
	userService,err := models.NewUserService(psqlInfo)
	defer userService.Close()
	if err != nil{
		panic(err)
	}
	userService.DestructiveReset()
	// user,err := userService.ByID(1)
	// fmt.Println(user)

}
