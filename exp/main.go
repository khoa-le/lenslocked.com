package main

import(
	"crypto/hmac"
	 "crypto/sha256"
	 "encoding/base64"
	 "fmt"

	 "lenslocked.com/hash"
)

func main(){
	toHash := []byte("vGcywVLvrT8jmlA9Gu798u3svYeaKNT1gyXcC7Xf_Uk=")
	h:= hmac.New(sha256.New, []byte("secret-hmac-key"))
	h.Write(toHash)

	b:=h.Sum(nil)
	fmt.Println(base64.URLEncoding.EncodeToString(b))

	hmac := hash.NewHMAC("my-secret-key")
	fmt.Println(hmac.Hash("this is my string to hash"))
}
