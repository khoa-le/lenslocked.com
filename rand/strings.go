package rand

import(
	"crypto/rand"
	"encoding/base64"
)

const RemeberTokenBytes =32

//Bytes help us generate n random bytes
func Bytes(n int) ([]byte,error){
	b :=make([]byte,n)
	_,err := rand.Read(b)
	if err != nil{
		return nil, err
	}
	return b, nil
}

func String(nBytes int) (string,error){
	b,err := Bytes(nBytes)
	if err != nil{
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func RememberToken() (string, error){
	return String(RemeberTokenBytes)
}