package models
import(
	"fmt"
	"testing"
	"time"
)

func testingUserService() (*UserService,error){
	const (
		host = "localhost"
		port = "5432"
		user = "khoa"
		dbname = "lenslocked_test"
	)

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
	host, port, user, dbname)

	us,err := NewUserService(psqlInfo)
	if err != nil{
		return nil,err
	}
	us.db.LogMode(true)
	//Clear the users table between tests
	us.DestructiveReset()
	return us, nil
}

func TestCreateUser(t *testing.T){
	us,err := testingUserService()
	if err != nil{
		t.Fatal(err)
	}
	user := User{
		Name: "Khanh Minh",
		Email: "khanhminh@gmail.com",
	}
	if err:= us.Create(&user); err !=nil{
		t.Fatal(err)
	}
	if user.ID == 0{
		t.Errorf("Expected >0. Recieved ID:%d",user.ID)
	}
	if time.Since(user.CreatedAt) >time.Duration(5*time.Second){
		t.Errorf("Expected Created At to be recent. Recieved %s", user.CreatedAt)
	}

}