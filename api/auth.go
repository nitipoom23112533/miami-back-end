package api
import(
	"github.com/golang-jwt/jwt/v5"

)
	

const jwtKey = "jwtKey"

type JWTCustomClaims struct {
	UID       string `json:"uid"`
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	IsAdmin   bool   `json:"isAdmin"`
	jwt.RegisteredClaims

}

func (c *JWTCustomClaims) Valid() error {
	return nil
}

func ParseJWTCustomClaims(x interface{}) *JWTCustomClaims {
	return x.(*jwt.Token).Claims.(*JWTCustomClaims)
}
func GetJwtKet() []byte  {
	return []byte(jwtKey)
}
