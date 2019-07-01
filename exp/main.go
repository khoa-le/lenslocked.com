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

	user := models.User{
		Name: "Khanh Minh 2",
		Email: "khanhminh3@gmail.com",
	}
	if err:= userService.Create(&user); err !=nil{
		panic(err)
	}

	// user.Email ="khanhminh2@gmail.com"
	// if err := userService.Update(&user); err != nil{
	// 	panic(err)
	// }

	userByEmail, err := userService.ByEmail("khanhminh3@gmail.com")
	if err !=nil{
		panic(err)
	}
	fmt.Println(userByEmail.Name)

	if err := userService.Delete(userByEmail.ID); err!=nil{
		panic(err)
	}

	// user,err := userService.ByID(1)
	// if err != nil{
	// 	panic(err)
	// }
	// fmt.Println(user)

}
