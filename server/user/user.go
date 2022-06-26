package user

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/RakaiSeto/simple-app-may/db"
	proto "github.com/RakaiSeto/simple-app-may/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

var dbconn *sql.DB
type RequestBody proto.RequestBody

func init() {
	dbconn = db.Db
}

func AllUser() (*proto.ResponseWrapper, error) {
	rows, err := dbconn.Query("SELECT id, uname, email, role FROM public.user ORDER BY id ASC")
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set"){
			var errString string = err.Error()
			return &proto.ResponseWrapper{
				Code: 404,
				Message: "not found",
				ResponseBody: &proto.ResponseBody{
					Error: &errString,
				},
			}, status.Error(codes.Code(5), "user not found")
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
		err := rows.Scan(&user.Id, &user.Uname, &user.Email, &user.Role)
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
	row := dbconn.QueryRow("SELECT id, uname, email, role FROM public.user where id=$1", id)

	user := proto.User{}
	err := row.Scan(&user.Id, &user.Uname, &user.Email, &user.Role)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set"){
			var errString string = err.Error()
			return &proto.ResponseWrapper{
				Code: 404,
				Message: "not found",
				ResponseBody: &proto.ResponseBody{
					Error: &errString,
				},
			}, status.Error(codes.Code(5), "user not found")
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

func AddUser(user *proto.User) (*proto.ResponseWrapper, error) {
	row := dbconn.QueryRow("SELECT id FROM public.user WHERE uname=$1", user.GetUname())
	var i int
	user.Role = "customer"
	err := row.Scan(&i)
	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result set") {
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
		}, status.Error(codes.Code(6), "user already exists, please change its uname")
	}
	
	_, err = dbconn.Exec("INSERT INTO public.user (uname, email, password, role, created, updated) VALUES ($1, $2, $3, $4, $5, $6)", user.GetUname(), user.GetEmail(), user.GetPassword(), "customer", time.Now().Unix(), time.Now().Unix())
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

	row = dbconn.QueryRow("SELECT id FROM public.user WHERE uname=$1", user.GetUname())
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

func UpdateUser(user *proto.User) (*proto.ResponseWrapper, error){
	QueryUser := proto.User{}
	row := dbconn.QueryRow("SELECT * from public.user where id = $1", user.Id)
	err := row.Scan(&QueryUser.Id, &QueryUser.Uname, &QueryUser.Email, &QueryUser.Password, &QueryUser.Role)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set"){
			var errString string = err.Error()
			return &proto.ResponseWrapper{
				Code: 404,
				Message: "not found",
				ResponseBody: &proto.ResponseBody{
					Error: &errString,
				},
			}, status.Error(codes.Code(5), "user not found")
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

	if QueryUser.GetPassword() != user.GetPassword() {
	var errString string = "wrong password for user"
		return &proto.ResponseWrapper{
			Code: 401,
			Message: "unauthorized",
			ResponseBody: &proto.ResponseBody{
				Error: &errString,
			},
		}, errors.New("wrong password for user")
	}

	if user.Uname != "" {QueryUser.Uname = user.Uname}
	if user.Email != "" {QueryUser.Email = user.Email}
	if *user.Password != "" {QueryUser.Password = user.Password}

	_, err = dbconn.Exec("UPDATE public.user SET uname=$2, email=$3, password=$4, updated=$5 WHERE id=$1", QueryUser.Id, QueryUser.Uname, QueryUser.Email, QueryUser.Password, time.Now().Unix())
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
			}, status.Error(codes.Code(5), "user not found")
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
			}, status.Error(codes.Code(5), "please input password")
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

	return &proto.ResponseWrapper{Code: 200, Message:"success", ResponseBody: &proto.ResponseBody{ResponseStatus: &proto.ResponseStatus{Response: "Success"}}}, nil
}

func CheckRequest(input string) (*proto.ResponseWrapper, error) {
	row := dbconn.QueryRow("SELECT reqbody, status FROM public.queue WHERE id = $1", input)
	var message string
	var byteSlice []byte
	if err := row.Scan(&byteSlice, &message); err != nil {
		fmt.Println(1)
		errstring := err.Error()
		return &proto.ResponseWrapper{Code: 500, Message: "Unknown", ResponseBody: &proto.ResponseBody{Error: &errstring}}, err
	}
	
	reqBody := new(proto.ResponseBody)
	err := protojson.Unmarshal(byteSlice, reqBody)
	if err != nil {
		fmt.Println(2)
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