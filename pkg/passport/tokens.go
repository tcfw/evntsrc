package passport

import (
	"crypto/rsa"
	fmt "fmt"
	"io/ioutil"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	users "github.com/tcfw/evntsrc/pkg/users/protos"
)

var (
	tlsKeyDir = "./"
)

//UserClaims takes in a user and applies the standard user claims
func UserClaims(user *users.User) map[string]interface{} {
	claims := make(map[string]interface{})

	claims["sub"] = user.Id
	claims["name"] = user.Name
	claims["email"] = user.Email
	claims["picture"] = user.Picture
	claims["groups"] = []string{}

	return claims
}

//GetKeyPrivate reads a private PEM formatted RSA cert
func GetKeyPrivate() (*rsa.PrivateKey, error) {
	dat, _ := ioutil.ReadFile(tlsKeyDir + "/priv.pem")
	key, er := jwt.ParseRSAPrivateKeyFromPEM(dat)
	if er != nil {
		fmt.Printf("Err: %s", er)
		return nil, er
	}

	return key, nil
}

//GetKeyPublic reads a public PEM formatted RSA cert
func GetKeyPublic() (*rsa.PublicKey, error) {
	dat, _ := ioutil.ReadFile(tlsKeyDir + "/pub.pem")
	key, er := jwt.ParseRSAPublicKeyFromPEM(dat)
	if er != nil {
		fmt.Printf("Err: %s", er)
		return nil, er
	}

	return key, nil
}

//MakeNewToken creates a new JWT token for the specific user
func MakeNewToken(extraClaims map[string]interface{}) (*string, *jwt.Token, error) {
	signer := jwt.New(jwt.SigningMethodRS256)

	//set claims
	claims := make(jwt.MapClaims)
	hostname, err := os.Hostname()
	if err != nil {
		claims["iss"] = "passport.gfc.io"
	} else {
		claims["iss"] = hostname + ".gfc.io"
	}
	claims["nbf"] = time.Now().Unix() - 1
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Hour * 168).Unix() // 1 week
	claims["jti"] = uuid.New().String()

	for claimKey, claimValue := range extraClaims {
		claims[claimKey] = claimValue
	}

	signer.Claims = claims

	key, err := GetKeyPrivate()
	if err != nil {
		return nil, nil, fmt.Errorf("Error extracting the key")
	}
	signer.Header["kid"] = "4b65488bcf182b9baead9e7c625763c104829912"

	tokenString, err := signer.SignedString(key)
	if err != nil {
		return nil, nil, err
	}

	return &tokenString, signer, nil
}

/*MakeTestToken makes a temporary JWT token which expires in 4
 *seconds specifically for testing purposes
 */
func MakeTestToken(user *users.User) (string, error) {
	signer := jwt.New(jwt.SigningMethodHS256)

	claims := make(jwt.MapClaims)
	claims["iss"] = "testing.gfc.io"
	claims["nbf"] = time.Now().Unix() - 1
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Second * 4).Unix()
	claims["jti"] = uuid.New().String()

	extraClaims := UserClaims(user)

	for claimKey, claimValue := range extraClaims {
		claims[claimKey] = claimValue
	}

	signer.Claims = claims
	return signer.SignedString([]byte("this is a super secret key, DO NOT USE FOR PRODUCTION"))
}

//MakeNewRefresh creates a new refresh token
func MakeNewRefresh() *string {
	randomString := RandomString(256)
	return &randomString
}
