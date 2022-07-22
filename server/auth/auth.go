package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/RakaiSeto/simple-app-may/db"
	proto "github.com/RakaiSeto/simple-app-may/service"
	"github.com/RakaiSeto/simple-app-may/helper"
	"github.com/antonholmquist/jason"
	redis "github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

var(
	redisconn *redis.Client
	ctx = context.TODO()
	dbconn *sql.DB
	githubOauthLogin *oauth2.Config
	githubState string
	googleOauthLogin *oauth2.Config
	googleState string
	facebookOauthLogin *oauth2.Config
	facebookState string
)

func init() {
	// loads values from .env into the system
    if err := godotenv.Load("auth/var.env"); err != nil {
        panic(err)
    }
	dbconn = db.Db
	redisconn = db.Rdb

	githubOauthLogin = &oauth2.Config{
		RedirectURL: "http://localhost:8080/login/github/callback",
		ClientID: os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Scopes: []string{"read:user"},
		Endpoint: github.Endpoint,
	}
	githubState = uuid.New().String()
	googleOauthLogin = &oauth2.Config{
		RedirectURL: "http://localhost:8080/login/google/callback",
		ClientID: os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint: google.Endpoint,
	}
	googleState = uuid.New().String()
	facebookOauthLogin = &oauth2.Config{
		RedirectURL: "http://localhost:8080/login/facebook/callback",
		ClientID: os.Getenv("FACEBOOK_CLIENT_ID"),
		ClientSecret: os.Getenv("FACEBOOK_CLIENT_SECRET"),
		Scopes: []string{"email", "public_profile"},
		Endpoint: facebook.Endpoint,
	}
	facebookState = uuid.New().String()
}

func Login(user *proto.User) (*proto.ResponseWrapper, error) {
	unameInsert := strings.ReplaceAll(user.GetUname(), " ", "")
	unameInsert = strings.ToLower(unameInsert)
	userRow := dbconn.QueryRow("SELECT id, password, role FROM public.user where uname=$1", unameInsert)
	var id int
	var password string
	var role string
	err := userRow.Scan(&id, &password, &role)
	if err != nil {
		if len(user.GetUname()) == 0 {
			var errString string = "please include the uname field in request"
			return &proto.ResponseWrapper{
				Code: 422,
				Message: "unprocessable entity",
				ResponseBody: &proto.ResponseBody{
					Error: &errString,
				},
			}, errors.New("please include uname")
		} else if strings.Contains(err.Error(), "no rows in result set"){
			var errString string = err.Error()
			return &proto.ResponseWrapper{
				Code: 404,
				Message: "not found",
				ResponseBody: &proto.ResponseBody{
					Error: &errString,
				},
			}, nil
		}
	}

	if password != user.GetPassword() {
		fmt.Println(password)
		fmt.Println(user.GetPassword())
		if user.GetPassword() == "" {
			var errString string = "please include password in request"
			return &proto.ResponseWrapper{
				Code: 422,
				Message: "unprocessable entity",
				ResponseBody: &proto.ResponseBody{
					Error: &errString,
				},
			}, errors.New("please include password in request")
		}
		var errString string = "wrong password for user"
		return &proto.ResponseWrapper{
			Code: 401,
			Message: "unauthorized",
			ResponseBody: &proto.ResponseBody{
				Error: &errString,
			},
		}, errors.New("wrong password for user")
	}

	if helper.IsJWTExist(ctx, user.GetUname(), redisconn){
		dbToken, err := helper.CheckJWT(ctx, user.GetUname(), redisconn)
		if err != nil {
			var errString string = err.Error()
			return &proto.ResponseWrapper{Code: 500, Message: "unknown error", ResponseBody: &proto.ResponseBody{Error: &errString}}, err
		}

		return &proto.ResponseWrapper{Code: 200, Message: "success", ResponseBody: &proto.ResponseBody{ResponseStatus: &proto.ResponseStatus{Response: "login success"}, String_: dbToken}}, nil
	}

	var tokenstring string
	if role == "admin" {
		tokenstring, err = helper.GenerateJWT(user.GetUname(), id, false, true)
		if err != nil {
			var errString string = err.Error()
			return &proto.ResponseWrapper{Code: 500, Message: "unknown error", ResponseBody: &proto.ResponseBody{Error: &errString}}, err
		}
	} else {	
		tokenstring, err = helper.GenerateJWT(user.GetUname(), id, false, false)
		if err != nil {
			var errString string = err.Error()
			return &proto.ResponseWrapper{Code: 500, Message: "unknown error", ResponseBody: &proto.ResponseBody{Error: &errString}}, err
		}
	}

	if err = inputJWT(user.GetUname(), tokenstring); err != nil {
		var errString string = err.Error()
		return &proto.ResponseWrapper{Code: 500, Message: "unknown error", ResponseBody: &proto.ResponseBody{Error: &errString}}, err
	}

	welcome := "login success, Welcome: " +user.GetUname()
	return &proto.ResponseWrapper{Code: 200, Message: "success", ResponseBody: &proto.ResponseBody{ResponseStatus: &proto.ResponseStatus{Response: welcome}, String_: &tokenstring}}, nil
}

func LoginGithub() (*proto.ResponseWrapper, error) {
	url := githubOauthLogin.AuthCodeURL(githubState)

	return &proto.ResponseWrapper{Code: 200, Message: "success", ResponseBody: &proto.ResponseBody{ResponseStatus: &proto.ResponseStatus{Response: url}}}, nil
}

func LoginGithubCallback(state string, code string) (*proto.ResponseWrapper, error) {
	if state != githubState {
		errString := "state string is different"
		return &proto.ResponseWrapper{Code: 401, Message:"unauthorized", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil
	}
	fmt.Println("CALLBACK")
	data, err := getGithubInfo(code)
	if err != nil {
		errString := err.Error()
		return &proto.ResponseWrapper{Code: 500, Message:"unknown", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil	
	}

	var x map[string]interface{}
	json.Unmarshal([]byte(data), &x)
	fmt.Printf("%v", x)
	fmt.Printf("%v, %t", x["login"], x["login"])

	resp, err := loginForOauth(x["login"].(string), "")
	if err != nil {
		errString := err.Error()
		return &proto.ResponseWrapper{Code: 500, Message:"unknown", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil
	}

	welcome := "Welcome: " + x["login"].(string)
	return &proto.ResponseWrapper{Code: 200, Message:"success", ResponseBody: &proto.ResponseBody{ResponseStatus: &proto.ResponseStatus{Response: welcome}, String_: &resp}}, nil
}

func LoginGoogle() (*proto.ResponseWrapper, error) {
	url := googleOauthLogin.AuthCodeURL(googleState)

	return &proto.ResponseWrapper{Code: 200, Message: "success", ResponseBody: &proto.ResponseBody{ResponseStatus: &proto.ResponseStatus{Response: url}}}, nil
}

func LoginGoogleCallback(state string, code string) (*proto.ResponseWrapper, error) {
	if state != googleState {
		errString := "state string is different"
		return &proto.ResponseWrapper{Code: 401, Message:"unauthorized", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil
	}
	fmt.Println("CALLBACK")
	data, err := getGoogleInfo(code)
	if err != nil {
		errString := err.Error()
		return &proto.ResponseWrapper{Code: 500, Message:"unknown", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil	
	}

	var x map[string]interface{}
	json.Unmarshal([]byte(data), &x)
	fmt.Printf("%v", x)

	resp, err := loginForOauth(x["name"].(string), x["email"].(string))
	if err != nil {
		errString := err.Error()
		return &proto.ResponseWrapper{Code: 500, Message:"unknown", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil
	}

	welcome := "Welcome: " + x["name"].(string)
	return &proto.ResponseWrapper{Code: 200, Message:"success", ResponseBody: &proto.ResponseBody{ResponseStatus: &proto.ResponseStatus{Response: welcome}, String_: &resp}}, nil
}

func LoginFacebook() (*proto.ResponseWrapper, error) {
	url := facebookOauthLogin.AuthCodeURL(facebookState)

	return &proto.ResponseWrapper{Code: 200, Message: "success", ResponseBody: &proto.ResponseBody{ResponseStatus: &proto.ResponseStatus{Response: url}}}, nil
}

func LoginFacebookCallback(state string, code string) (*proto.ResponseWrapper, error) {
	if state != facebookState {
		errString := "state string is different"
		return &proto.ResponseWrapper{Code: 401, Message:"unauthorized", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil
	}
	fmt.Println("CALLBACK")
	data, err := getFacebookInfo(code)
	if err != nil {
		errString := err.Error()
		return &proto.ResponseWrapper{Code: 200, Message:"success", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil	
	}

	bs, err := data.Marshal()
	if err != nil {
		panic(err)
	}
	var x map[string]interface{}
	json.Unmarshal(bs, &x)
	fmt.Printf("%v", x)

	resp, err := loginForOauth(x["name"].(string), x["email"].(string))
	if err != nil {
		errString := err.Error()
		return &proto.ResponseWrapper{Code: 500, Message:"unknown", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil
	}

	welcome := "Welcome: " + x["name"].(string)
	return &proto.ResponseWrapper{Code: 200, Message:"success", ResponseBody: &proto.ResponseBody{ResponseStatus: &proto.ResponseStatus{Response: welcome}, String_: &resp}}, nil
}

func Logout(tokenString string) (*proto.ResponseWrapper, error){
	creden, err := helper.ParseJWT(ctx, tokenString)
	if err != nil{
		var errString string = err.Error()
		return &proto.ResponseWrapper{Code: 500, Message: "unknown error", ResponseBody: &proto.ResponseBody{Error: &errString}}, err
	}

	err = helper.DeleteJWT(ctx, creden["user"].(string), redisconn)
	if err != nil {
		var errString string = err.Error()
		return &proto.ResponseWrapper{Code: 500, Message: "unknown error", ResponseBody: &proto.ResponseBody{Error: &errString}}, err
	}

	return &proto.ResponseWrapper{Code: 200, Message: "success", ResponseBody: &proto.ResponseBody{ResponseStatus: &proto.ResponseStatus{Response: "success"}}}, nil
}

func inputJWT(uname string, token string) (error) {
	err := redisconn.HSet(ctx, "jwtdb", uname, token).Err()
	if err != nil {
		return err
	}
	return nil
}

func getGithubInfo(code string) (string, error) {
	token, err := githubOauthLogin.Exchange(context.TODO(), code)
    if err != nil {
		return "", fmt.Errorf("code exchange failed: %s", err.Error())
	}

    req, reqerr := http.NewRequest(
        "GET",
        "https://api.github.com/user",
        nil,
    )
	if reqerr != nil {
        panic(reqerr)
    }
	authorizationHeaderValue := fmt.Sprintf("token %s", token.AccessToken)
    req.Header.Set("Authorization", authorizationHeaderValue)

    // Get the response
	var httpclient http.Client
    resp, resperr := httpclient.Do(req)
    if resperr != nil {
        panic(resperr)
    }
	defer resp.Body.Close()

    // Response body converted to stringified JSON
	
	fmt.Println(resp.Status)
	fmt.Println(resp.Body)
    respbody, _ := ioutil.ReadAll(resp.Body)

	return string(respbody), nil
}

func getGoogleInfo(code string) (string, error) {
	token, err := googleOauthLogin.Exchange(context.TODO(), code)
    if err != nil {
		return "", fmt.Errorf("code exchange failed: %s", err.Error())
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
    if err != nil {
        panic(err)
    }
	defer resp.Body.Close()

    // Response body converted to stringified JSON
	
	fmt.Println(resp.Status)
	fmt.Println(resp.Body)
    respbody, _ := ioutil.ReadAll(resp.Body)

	return string(respbody), nil
}

func getFacebookInfo(code string) (*jason.Object, error) {
	token, err := facebookOauthLogin.Exchange(context.TODO(), code)
    if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	resp, err := http.Get("https://graph.facebook.com/me?fields=email,name&access_token=" + token.AccessToken)
    if err != nil {
        panic(err)
    }
	defer resp.Body.Close()

	bodyBuffer := make([]byte, 5000)
 	var str string

 	count, err := resp.Body.Read(bodyBuffer)

 	for ; count > 0; count, err = resp.Body.Read(bodyBuffer) {

 		if err != nil {
			panic(err)
 		}

 		str += string(bodyBuffer[:count])
 	}

    // Response body converted to stringified JSON
	user, _ := jason.NewObjectFromBytes([]byte(str))

	return user, nil
}

func loginForOauth(user string, email string) (string, error) {
	userRow := dbconn.QueryRow("SELECT id, uname, password FROM public.user WHERE email=$1", email)
	var id int
	var dbname string
	var password string
	err := userRow.Scan(&id, &dbname, &password)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set"){
			userRow := dbconn.QueryRow("SELECT id, email, password FROM public.user WHERE uname=$1", user)
			var id int
			var dbemail string
			var password string
			err := userRow.Scan(&id, &dbemail, &password)
			if err != nil {
				if strings.Contains(err.Error(), "no rows in result set"){
					_, err = dbconn.Exec("INSERT INTO public.user (uname, email, created_at, updated_at) VALUES ($1, $2, $3, $4)", user, email, time.Now().Unix(), time.Now().Unix())
					if err != nil {
						return err.Error(), err
					}
					userRow := dbconn.QueryRow("SELECT id WHERE uname=$1", user)
					err := userRow.Scan(&id)
					if err != nil {
						return err.Error(), err
					}
					var id int
					respsonse, err := helper.GenerateJWT(user, id, true, false)
					if err != nil {
						return err.Error(), err
					}
					inputJWT(user, respsonse)
					if err != nil {
						return err.Error(), err
					}
					return respsonse, nil
				}
			}
			respsonse, err := helper.GenerateJWT(user, id, true, false)
			if err != nil {
				return err.Error(), err
			}
			inputJWT(user, respsonse)
			if err != nil {
				return err.Error(), err
			}
			return respsonse, nil
		}
	}
	respsonse, err := helper.GenerateJWT(dbname, id, true, false)
	if err != nil {
		return err.Error(), err
	}
	inputJWT(dbname, respsonse)
	if err != nil {
		return err.Error(), err
	}
	return respsonse, nil
}