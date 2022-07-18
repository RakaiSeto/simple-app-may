package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/RakaiSeto/simple-app-may/db"
	proto "github.com/RakaiSeto/simple-app-may/service"
	"google.golang.org/protobuf/encoding/protojson"
)

var dbconn *sql.DB
type RequestBody proto.RequestBody

var funcCtx = context.TODO()

func init() {
	dbconn = db.Db
}

func AllUser() (*proto.ResponseWrapper, error) {
	rows, err := dbconn.Query("SELECT id, uname, email FROM public.user ORDER BY id ASC")
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set"){
			var errString string = err.Error()
			return &proto.ResponseWrapper{
				Code: 404,
				Message: "not found",
				ResponseBody: &proto.ResponseBody{
					Error: &errString,
				},
			}, nil
		}
		var errString string = err.Error()
		return &proto.ResponseWrapper{
			Code: 500,
			Message: "unknown error",
			ResponseBody: &proto.ResponseBody{
				Error: &errString,
			},
		}, err 
	}
	defer rows.Close()

	users := make([]*proto.User, 0)
	for rows.Next() {
		var user proto.User
		err := rows.Scan(&user.Id, &user.Uname, &user.Email)
		if err != nil {
			var errString string = err.Error()
			return &proto.ResponseWrapper{
				Code: 500,
				Message: "unknown error",
				ResponseBody: &proto.ResponseBody{
					Error: &errString,
				},
			}, err
		}

		users = append(users, &user)
	}

	returned := &proto.ResponseWrapper{Code: 200, Message: "success", ResponseBody: &proto.ResponseBody{
		Users: &proto.Users{
			User: users,
		},
	}}

	return returned, nil
}

func OneUser(id int) (*proto.ResponseWrapper, error) {
	row := dbconn.QueryRow("SELECT id, uname, email FROM public.user where id=$1", id)

	user := proto.User{}
	err := row.Scan(&user.Id, &user.Uname, &user.Email)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set"){
			var errString string = err.Error()
			return &proto.ResponseWrapper{
				Code: 404,
				Message: "not found",
				ResponseBody: &proto.ResponseBody{
					Error: &errString,
				},
			}, nil
		}
		var errString string = err.Error()
		return &proto.ResponseWrapper{
			Code: 500,
			Message: "unknown error",
			ResponseBody: &proto.ResponseBody{
				Error: &errString,
			},
		}, err
	}

	returned := &proto.ResponseWrapper{Code: 200, Message: "success", ResponseBody: &proto.ResponseBody{
		User: &user,
	}}

	return returned, nil
}

func MyUser(token string) (*proto.ResponseWrapper, error) {
	creden, err := proto.ParseJWT(funcCtx, token)
	if err != nil {
		var errString string = err.Error()
		return &proto.ResponseWrapper{Code: 401, Message: "unauthorized", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil
	}

	row := dbconn.QueryRow("SELECT * FROM public.user where id=$1", creden["userid"].(float64))

	user := proto.User{}

	var password string
	var created int64
	var updated int64
	err = row.Scan(&user.Id, &user.Uname, &user.Email, &password, &user.Role, &created, &updated)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set"){
			var errString string = err.Error()
			return &proto.ResponseWrapper{
				Code: 404,
				Message: "not found",
				ResponseBody: &proto.ResponseBody{
					Error: &errString,
				},
			}, nil
		}
		var errString string = err.Error()
		return &proto.ResponseWrapper{
			Code: 500,
			Message: "unknown error",
			ResponseBody: &proto.ResponseBody{
				Error: &errString,
			},
		}, err
	}

	
	createdString := time.Unix(created, 0).In(proto.WIB_TIME).Format(proto.TIME_LAYOUT_ALL)
	updatedString := time.Unix(updated, 0).In(proto.WIB_TIME).Format(proto.TIME_LAYOUT_ALL)

	user.CreatedAt = &createdString
	user.UpdatedAt = &updatedString

	returned := &proto.ResponseWrapper{Code: 200, Message: "success", ResponseBody: &proto.ResponseBody{
		User: &user,
	}}

	return returned, nil
}

func AddUser(user *proto.User) (*proto.ResponseWrapper, error) {
	row := dbconn.QueryRow("SELECT id FROM public.user WHERE uname=$1", user.GetUname())
	var i int
	*user.Role = "customer"
	err := row.Scan(&i)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
		} else {
			var errString string = err.Error()
			return &proto.ResponseWrapper{
				Code: 500,
				Message: "unknown error",
				ResponseBody: &proto.ResponseBody{
					Error: &errString,
				},
				}, err
		}
	}
	if i != 0 {
		var errString string = "user already exists"
		return &proto.ResponseWrapper{
			Code: 409,
			Message: "conflict",
			ResponseBody: &proto.ResponseBody{
				Error: &errString,
			},
		}, err
	}
	
	unameInsert := strings.ReplaceAll(user.GetUname(), " ", "")
	unameInsert = strings.ToLower(unameInsert)

	_, err = dbconn.Exec("INSERT INTO public.user (uname, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)", unameInsert, user.GetEmail(), user.GetPassword(), time.Now().Unix(), time.Now().Unix())
	if err != nil {
		var errString string = err.Error()
		return &proto.ResponseWrapper{
			Code: 500,
			Message: "unknown error",
			ResponseBody: &proto.ResponseBody{
				Error: &errString,
			},
			}, err
	}

	row = dbconn.QueryRow("SELECT id FROM public.user WHERE uname=$1", unameInsert)
	err = row.Scan(&user.Id)
	if err != nil {
		var errString string = err.Error()
		return &proto.ResponseWrapper{
			Code: 500,
			Message: "unknown error",
			ResponseBody: &proto.ResponseBody{
				Error: &errString,
			},
		}, err
	}

	resp := proto.AddUserStatus{Response: "success", User: user}

	return &proto.ResponseWrapper{
		Code: 200,
		Message: "success", 
		ResponseBody: &proto.ResponseBody{
			AddUserStatus: &resp,
		},
	}, nil
}

func UpdateUser(input *proto.RequestBody) (*proto.ResponseWrapper, error){
	creden, err := proto.ParseJWT(funcCtx, input.GetString_())
	if err != nil {
		var errString string = err.Error()
		return &proto.ResponseWrapper{Code: 401, Message: "unauthorized", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil
	}

	QueryUser := proto.User{}
	row := dbconn.QueryRow("SELECT id, uname, email, password from public.user where id = $1", creden["userid"].(float64))
	err = row.Scan(&QueryUser.Id, &QueryUser.Uname, &QueryUser.Email, &QueryUser.Password)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set"){
			var errString string = err.Error()
			return &proto.ResponseWrapper{
				Code: 404,
				Message: "not found",
				ResponseBody: &proto.ResponseBody{
					Error: &errString,
				},
			}, err
		}
		var errString string = err.Error()
		return &proto.ResponseWrapper{
			Code: 500,
			Message: "unknown error",
			ResponseBody: &proto.ResponseBody{
				Error: &errString,
			},
		}, err
	}

	if input.User.Uname != "" {QueryUser.Uname = input.User.Uname}
	QueryUser.Uname = strings.ReplaceAll(QueryUser.Uname, " ", "")
	QueryUser.Uname = strings.ToLower(QueryUser.Uname)
	if input.User.Email != "" {QueryUser.Email = input.User.Email}
	if *input.User.Password != "" {QueryUser.Password = input.User.Password}

	_, err = dbconn.Exec("UPDATE public.user SET uname=$2, email=$3, password=$4, updated_at=$5 WHERE id=$1", QueryUser.Id, QueryUser.Uname, QueryUser.Email, QueryUser.Password, time.Now().Unix())
	if err != nil {
		var errString string = err.Error()
		return &proto.ResponseWrapper{
			Code: 500,
			Message: "unknown error",
			ResponseBody: &proto.ResponseBody{
				Error: &errString,
			},
		}, err
	}
	return &proto.ResponseWrapper{Code: 200, Message:"success", ResponseBody: &proto.ResponseBody{ResponseStatus: &proto.ResponseStatus{Response: "Success"}}}, nil
}

func DeleteUser(inputUser *proto.User) (*proto.ResponseWrapper, error) {
	row := dbconn.QueryRow("SELECT uname, password FROM public.user where id=$1", inputUser.Id)

	var name string
	var password string
	
	err := row.Scan(&name, &password)
	if err != nil {
		var errString string = err.Error()
		return &proto.ResponseWrapper{
			Code: 500,
			Message: "unknown error",
			ResponseBody: &proto.ResponseBody{
				Error: &errString,
			},
		}, err
	}

	if name == ""{
		var errString string = "user not found"
			return &proto.ResponseWrapper{
				Code: 404,
				Message: "not found",
				ResponseBody: &proto.ResponseBody{
					Error: &errString,
				},
			}, err
	}
	
	if inputUser.GetPassword() != password {
		if inputUser.GetPassword() == "" {
			var errString string = "please input password"
			return &proto.ResponseWrapper{
				Code: 422,
				Message: "unprocessable entity",
				ResponseBody: &proto.ResponseBody{
					Error: &errString,
				},
			}, errors.New("please input password")
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

	_, err = dbconn.Exec("DELETE FROM public.user WHERE id=$1", inputUser.Id)
	if err != nil {
		var errString string = err.Error()
		return &proto.ResponseWrapper{
			Code: 500,
			Message: "unknown error",
			ResponseBody: &proto.ResponseBody{
				Error: &errString,
			},
		}, err
	}

	proto.DeleteJWT(context.TODO(), name, db.Rdb)

	return &proto.ResponseWrapper{Code: 200, Message:"success", ResponseBody: &proto.ResponseBody{ResponseStatus: &proto.ResponseStatus{Response: "Success"}}}, nil
}

func AdminTopup(input *proto.AdminTopup) (*proto.ResponseWrapper, error) {
	row := dbconn.QueryRow("SELECT uname, wallet FROM public WHERE id=$1", input.Userid)

	var uname string
	var wallet float32

	err := row.Scan(&uname, &wallet)
	if err != nil {
		var errString string = err.Error()
		return &proto.ResponseWrapper{
			Code: 500,
			Message: "unknown error",
			ResponseBody: &proto.ResponseBody{
				Error: &errString,
			},
		}, nil
	}

	if uname != input.Username {
		var errString string = "the id have different username"
		return &proto.ResponseWrapper{
			Code: 404,
			Message: "not found",
			ResponseBody: &proto.ResponseBody{
				Error: &errString,
			},
		}, nil
	}

	wallet /= 100
	wallet += input.Amount
	wallet *= 100

	_, err = dbconn.Exec("UPDATE public.user SET wallet=$1, updated_at=$2 WHERE id=$3", wallet, time.Now().Unix(), input.Userid)
	if err != nil {
		var errString string = err.Error()
		return &proto.ResponseWrapper{
			Code: 500,
			Message: "unknown error",
			ResponseBody: &proto.ResponseBody{
				Error: &errString,
			},
		}, err
	}

	return &proto.ResponseWrapper{Code: 200, Message:"success", ResponseBody: &proto.ResponseBody{ResponseStatus: &proto.ResponseStatus{Response: "Success"}}}, nil
}

func CheckRequest(input string) (*proto.ResponseWrapper, error) {
	row := dbconn.QueryRow("SELECT reqbody, status FROM public.queue WHERE id = $1", input)
	var message string
	var byteSlice []byte
	if err := row.Scan(&byteSlice, &message); err != nil {
		errstring := err.Error()
		return &proto.ResponseWrapper{Code: 500, Message: "Unknown", ResponseBody: &proto.ResponseBody{Error: &errstring}}, err
	}
	
	reqBody := new(proto.ResponseBody)
	err := protojson.Unmarshal(byteSlice, reqBody)
	if err != nil {
		errstring := err.Error()
		return &proto.ResponseWrapper{Code: 500, Message: "Unknown", ResponseBody: &proto.ResponseBody{Error: &errstring}}, err
	}

	return &proto.ResponseWrapper{Code: 200, Message: "success", ResponseBody: reqBody}, nil
}

func (r *RequestBody) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok{
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &r)
}

func GetUser(user *proto.User, id int) (*proto.User, error) {
	row := dbconn.QueryRow("SELECT * FROM public.user WHERE id = $1", id)
	var created int
	var updated int
	err := row.Scan(&user.Id, &user.Uname, &user.Email, &user.Password, &user.Role, &created, &updated)
	if err != nil {
		return nil, err
	}

	return user, nil
}