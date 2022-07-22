package helper

import (
	context "context"
	"errors"
	"fmt"
	"time"

	"github.com/RakaiSeto/simple-app-may/db"
	"github.com/RakaiSeto/simple-app-may/service"
	"github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis/v9"
	jwt "github.com/golang-jwt/jwt"
)

var funcCtx = context.TODO()

var	signing_string = []byte("supersecret")

type Credentials struct {
	User string
	Password string
	Expired int64
}

func GenerateJWT(user string, userid int, oauth bool, admin bool) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["oauth"] = oauth
	claims["admin"] = admin
	claims["issuer"] = "RakaiApp"
	claims["user"] = user
	claims["userid"] = userid
	claims["expired"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString(signing_string)

	if err != nil {
		panic(err)
	}

	fmt.Println(tokenString)
	return tokenString, nil
}

func IsJWTExist(ctx context.Context, uname string, redisConn *redis.Client) bool {
	query, err := redisConn.HGet(ctx, "jwtdb", uname).Result()
	if ((err != nil) || (query == ""))  {
		if err != nil {
			fmt.Println(err.Error())
		}
		return false
	}
	return true
}


func CheckJWT(ctx context.Context, uname string, redisConn *redis.Client) (*string, error) {
	query, err := redisConn.HGet(ctx, "jwtdb", uname).Result()
	if err != nil {
		return nil, err
	}
	return &query, nil
}

func DeleteJWT(ctx context.Context, uname string, redisConn *redis.Client) (error) {
	bool_ := IsJWTExist(ctx, uname, redisConn)
	if !bool_{
		return errors.New("not logged in yet")
	}

	err := redisConn.HDel(ctx, "jwtdb", uname).Err()
	if err != nil {
		return err
	}
	return nil
}

func ParseJWT(ctx context.Context, token string) (jwt.MapClaims, error){
	mapToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return signing_string, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := mapToken.Claims.(jwt.MapClaims); ok && mapToken.Valid {
		return claims, nil
	}
	return nil, err
}

func ValidateJWT(ctx *gin.Context) *service.ResponseWrapper {
	cookie, err := ctx.Cookie("token")
	if err != nil {
		var errString string = "not login yet"
		return &service.ResponseWrapper{Code: 401, Message: "unauthorized", ResponseBody: &service.ResponseBody{Error: &errString}}
	}

	creden, err := ParseJWT(funcCtx, cookie)
	if err != nil {
		var errString string = "not login yet"
		return &service.ResponseWrapper{Code: 401, Message: "unauthorized", ResponseBody: &service.ResponseBody{Error: &errString}}
	}

	fmt.Println(creden)

	if creden["expired"].(float64) < float64(time.Now().Unix()) {
		var errString string = "please re-login"
		DeleteJWT(ctx, creden["user"].(string), db.Rdb)
		ctx.SetCookie("token", "", -1, "/", "localhost", false, true)
		return &service.ResponseWrapper{Code: 401, Message: "unauthorized", ResponseBody: &service.ResponseBody{Error: &errString}}
	}

	if !IsJWTExist(funcCtx, creden["user"].(string), db.Rdb){
		var errString string = "not login yet"
		return &service.ResponseWrapper{Code: 401, Message: "unauthorized", ResponseBody: &service.ResponseBody{Error: &errString}}
	}

	tokenString, err := CheckJWT(funcCtx, creden["user"].(string), db.Rdb)
	if err != nil {
		var errString string = err.Error()
		return &service.ResponseWrapper{Code: 401, Message: "unauthorized", ResponseBody: &service.ResponseBody{Error: &errString}}
	}

	if *tokenString != cookie {
		var errString string = "please re-login"
		return &service.ResponseWrapper{Code: 401, Message: "unauthorized", ResponseBody: &service.ResponseBody{Error: &errString}}
	}

	ctx.SetCookie("token", cookie, 30*60, "/", "localhost", false, true)

	return nil
}