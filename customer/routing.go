package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	proto "github.com/RakaiSeto/simple-app-may/service"
	"github.com/gin-gonic/gin"
)


func Tes(ctx *gin.Context) {
	reqBody := &proto.EmptyStruct{}
	req := &proto.RequestWrapper{Method: getMethod(ctx), Url: getURL(ctx), RequestBody: &proto.RequestBody{EmptyStruct: reqBody}}
	if response, _ := Client.Tes(ctx, req); response.ResponseBody.GetError() == "" {
		// Prettifying the json
		fmt.Println(response.ResponseBody.ResponseStatus.GetResponse())
		var prettyJSON bytes.Buffer
		// json.indent is a library utility function to prettify JSON indentation
		parserr := json.Indent(&prettyJSON, []byte(response.ResponseBody.ResponseStatus.GetResponse()), "", "\t")
		if parserr != nil {
			panic(parserr)
		}
		ctx.IndentedJSON(http.StatusOK, string(prettyJSON.Bytes()))
		return
	} else {
		errorHandler(ctx, 1, response)
        return
	}
}

func Login(ctx *gin.Context) {
	var user proto.User

	if err := ctx.BindJSON(&user); err != nil {
		var errString string = err.Error()
		ctx.IndentedJSON(http.StatusBadRequest, &proto.ResponseWrapper{Code: 400, Message: "bad request", ResponseBody: &proto.ResponseBody{Error: &errString}})
	}

	req := &proto.RequestWrapper{Method: getMethod(ctx), Url: getURL(ctx), RequestBody: &proto.RequestBody{User: &user}}

	if response, _ := Client.Login(ctx, req); response.ResponseBody.GetError() != "" {
		errorHandler(ctx, 1, response)
		return
	} else {
		ctx.SetCookie("token", response.ResponseBody.GetString_(), 30*60, "/", "localhost", false, true)
		response.ResponseBody.String_ = nil
		ctx.IndentedJSON(http.StatusOK, response)
	}
}

func LoginGithub(ctx *gin.Context) {
	cookie, _ := ctx.Cookie("token")
	if cookie != "" {
		var errString string = "please log out first"
		errorHandler(ctx, 1, &proto.ResponseWrapper{Code: 400, Message: "bad request", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
	}

	reqBody := &proto.EmptyStruct{}
	req := &proto.RequestWrapper{Method: getMethod(ctx), Url: getURL(ctx), RequestBody: &proto.RequestBody{EmptyStruct: reqBody}}
	if response, _ := Client.LoginGithub(ctx, req); response.ResponseBody.GetError() == "" {
		ctx.Redirect(http.StatusTemporaryRedirect, response.ResponseBody.ResponseStatus.GetResponse())
		return
	} else {
		errorHandler(ctx, 1, response)
        return
	}
}

func LoginGithubCallback(ctx *gin.Context) {
	state := ctx.Query("state")
	code := ctx.Query("code")
	req := &proto.RequestWrapper{Method: getMethod(ctx), Url: getURL(ctx), RequestBody: &proto.RequestBody{OauthCallback: &proto.OauthCallback{State: state, Code: code}}}
	if response, _ := Client.LoginGithubCallback(ctx, req); response.ResponseBody.GetError() == "" {
		ctx.IndentedJSON(http.StatusOK, response)
		return
	} else {
		errorHandler(ctx, 1, response)
        return
	}
}

func LoginGoogle(ctx *gin.Context) {
	cookie, _ := ctx.Cookie("token")
	if cookie != "" {
		var errString string = "please log out first"
		errorHandler(ctx, 1, &proto.ResponseWrapper{Code: 400, Message: "bad request", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
	}

	reqBody := &proto.EmptyStruct{}
	req := &proto.RequestWrapper{Method: getMethod(ctx), Url: getURL(ctx), RequestBody: &proto.RequestBody{EmptyStruct: reqBody}}
	if response, _ := Client.LoginGoogle(ctx, req); response.ResponseBody.GetError() == "" {
		ctx.Redirect(http.StatusTemporaryRedirect, response.ResponseBody.ResponseStatus.GetResponse())
		return
	} else {
		errorHandler(ctx, 1, response)
        return
	}
}

func LoginGoogleCallback(ctx *gin.Context) {
	state := ctx.Query("state")
	code := ctx.Query("code")
	req := &proto.RequestWrapper{Method: getMethod(ctx), Url: getURL(ctx), RequestBody: &proto.RequestBody{OauthCallback: &proto.OauthCallback{State: state, Code: code}}}
	if response, _ := Client.LoginGoogleCallback(ctx, req); response.ResponseBody.GetError() == "" {
		ctx.IndentedJSON(http.StatusOK, response)
		return
	} else {
		errorHandler(ctx, 1, response)
        return
	}
}

func LoginFacebook(ctx *gin.Context) {
	cookie, _ := ctx.Cookie("token")
	if cookie != "" {
		var errString string = "please log out first"
		errorHandler(ctx, 1, &proto.ResponseWrapper{Code: 400, Message: "bad request", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
	}

	reqBody := &proto.EmptyStruct{}
	req := &proto.RequestWrapper{Method: getMethod(ctx), Url: getURL(ctx), RequestBody: &proto.RequestBody{EmptyStruct: reqBody}}
	if response, _ := Client.LoginFacebook(ctx, req); response.ResponseBody.GetError() == "" {
		ctx.Redirect(http.StatusTemporaryRedirect, response.ResponseBody.ResponseStatus.GetResponse())
		return
	} else {
		errorHandler(ctx, 1, response)
        return
	}
}

func LoginFacebookCallback(ctx *gin.Context) {
	state := ctx.Query("state")
	code := ctx.Query("code")
	req := &proto.RequestWrapper{Method: getMethod(ctx), Url: getURL(ctx), RequestBody: &proto.RequestBody{OauthCallback: &proto.OauthCallback{State: state, Code: code}}}
	if response, _ := Client.LoginFacebookCallback(ctx, req); response.ResponseBody.GetError() == "" {
		ctx.IndentedJSON(http.StatusOK, response)
		return
	} else {
		errorHandler(ctx, 1, response)
        return
	}
}

func Logout(ctx *gin.Context) {
	cookie, err := ctx.Cookie("token")
	if err != nil {
		var errString string = "not logged in yet"
		errorHandler(ctx, 1, &proto.ResponseWrapper{Code: 422, Message: "unprocessable entity", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
	}

	req := &proto.RequestWrapper{Method: getMethod(ctx), Url: getURL(ctx), RequestBody: &proto.RequestBody{String_: &cookie}}

	if response, _ := Client.Logout(ctx, req); response.ResponseBody.GetError() != "" {
		errorHandler(ctx, 1, response)
		return
	} else {
		ctx.SetCookie("token", "", -1, "/", "localhost", false, true)
		ctx.IndentedJSON(http.StatusOK, response)
	}
}

func CheckError(ctx *gin.Context) {
	id := ctx.Param("id")
	req := &proto.RequestWrapper{Method: getMethod(ctx), Url: getURL(ctx), RequestBody: &proto.RequestBody{String_: &id}}
	if response, _ := Client.ReqError(ctx, req); response.ResponseBody.GetError() == "" {
		ctx.IndentedJSON(http.StatusOK, response)
		return
	} else {
		errorHandler(ctx, 1, response)
	}
}

func AllUser(ctx *gin.Context) {
	jwterr := proto.ValidateJWT(ctx)
	if jwterr != nil {
		errorHandler(ctx, 1, jwterr)
		return
	}

	reqBody := &proto.EmptyStruct{}
	req := &proto.RequestWrapper{Method: getMethod(ctx), Url: getURL(ctx), RequestBody: &proto.RequestBody{EmptyStruct: reqBody}}
	if response, _ := Client.AllUser(ctx, req); response.ResponseBody.GetError() == "" {
		ctx.IndentedJSON(http.StatusOK, response)
		return
	} else {	
		errorHandler(ctx,  1, response)
		return
	}
}

func OneUser(ctx *gin.Context) {
	jwterr := proto.ValidateJWT(ctx)
	if jwterr != nil {
		errorHandler(ctx, 1, jwterr)
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		var errString string = err.Error()
		ctx.IndentedJSON(http.StatusBadRequest, &proto.ResponseWrapper{Code: 400, Message: "bad request", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
	}

	reqBody := &proto.Id{Id: int64(id)}
	req := &proto.RequestWrapper{Method: getMethod(ctx), Url: getURL(ctx), RequestBody: &proto.RequestBody{Id: reqBody}}
	if response, _ := Client.OneUser(ctx, req); response.ResponseBody.GetError() != "" {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		errorHandler(ctx, 1, response)
	}
}

func PostUser(ctx *gin.Context) {
	var user proto.User

	if err := ctx.BindJSON(&user); err != nil {
		var errString string = err.Error()
		ctx.IndentedJSON(http.StatusInternalServerError, &proto.ResponseWrapper{Code: 500, Message: "unknown error", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
	}

	req := &proto.RequestWrapper{Method: getMethod(ctx), Url: getURL(ctx), RequestBody: &proto.RequestBody{User: &user}}
	if response, _ := Produce(req); response.ResponseBody.GetError() != "" {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		errorHandler(ctx, 1, response)
	}
}

func PatchUser(ctx *gin.Context) {
	jwterr := proto.ValidateJWT(ctx)
	if jwterr != nil {
		errorHandler(ctx, 1, jwterr)
		return
	}

	var user proto.User

	if err := ctx.BindJSON(&user); err != nil {
		var errString string = err.Error()
		ctx.IndentedJSON(http.StatusInternalServerError, &proto.ResponseWrapper{Code: 500, Message: "unknown error", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
    }

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		var errString string = err.Error()
		ctx.IndentedJSON(http.StatusBadRequest, &proto.ResponseWrapper{Code: 400, Message: "bad request", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
	}

	user.Id = int64(id)

	req := &proto.RequestWrapper{Method: getMethod(ctx), Url: getURL(ctx), RequestBody: &proto.RequestBody{User: &user}}
	if response, _ := Produce(req); response.ResponseBody.GetError() != "" {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		errorHandler(ctx, 1, response)
	}
}

func DeleteUser(ctx *gin.Context) {
	jwterr := proto.ValidateJWT(ctx)
	if jwterr != nil {
		errorHandler(ctx, 1, jwterr)
		return
	}

	var user proto.User

	if err := ctx.BindJSON(&user); err != nil {
		var errString string = err.Error()
		ctx.IndentedJSON(http.StatusInternalServerError, &proto.ResponseWrapper{Code: 500, Message: "unknown error", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
    }

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		var errString string = err.Error()
		ctx.IndentedJSON(http.StatusBadRequest, &proto.ResponseWrapper{Code: 400, Message: "bad request", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
	}

	user.Id = id

	req := &proto.RequestWrapper{Method: getMethod(ctx), Url: getURL(ctx), RequestBody: &proto.RequestBody{User: &user}}
	if response, _ := Produce(req); response.ResponseBody.GetError() != "" {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		errorHandler(ctx, 1, response)
	}
}

func AllProduct(ctx *gin.Context) {
	jwterr := proto.ValidateJWT(ctx)
	if jwterr != nil {
		errorHandler(ctx, 2, jwterr)
		return
	}

	reqbody := &proto.EmptyStruct{}
	req := &proto.RequestWrapper{Method: getMethod(ctx), Url: getURL(ctx), RequestBody: &proto.RequestBody{EmptyStruct: reqbody}}
	if response, err := Produce(req); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		errorHandler(ctx, 2, response)
	}
}

func OneProduct(ctx *gin.Context) {
	jwterr := proto.ValidateJWT(ctx)
	if jwterr != nil {
		errorHandler(ctx, 2, jwterr)
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		var errString string = err.Error()
		ctx.IndentedJSON(http.StatusBadRequest, &proto.ResponseWrapper{Code: 400, Message: "bad request", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
	}

	reqbody := &proto.Id{Id: id}
	req := &proto.RequestWrapper{Method: getMethod(ctx), Url: getURL(ctx), RequestBody: &proto.RequestBody{Id: reqbody}}
	if response, err := Produce(req); err == nil {
		ctx.IndentedJSON(http.StatusInternalServerError, response)
	} else {	
		errorHandler(ctx, 2, response)
	}
}

func AllOrder(ctx *gin.Context) {
	jwterr := proto.ValidateJWT(ctx)
	if jwterr != nil {
		errorHandler(ctx, 3, jwterr)
	}

	var user proto.User

	if err := ctx.BindJSON(&user); err != nil {
		var errString string = err.Error()
		ctx.IndentedJSON(http.StatusBadRequest, &proto.ResponseWrapper{Code: 400, Message: "bad request", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
    }

	req := &proto.RequestWrapper{Method: getMethod(ctx), Url: getURL(ctx), RequestBody: &proto.RequestBody{User: &user}}
	if response, err := Produce(req); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {
		errorHandler(ctx, 3, response)
	}
}

func OneOrder(ctx *gin.Context) {
	jwterr := proto.ValidateJWT(ctx)
	if jwterr != nil {
		errorHandler(ctx, 3, jwterr)
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		var errString string = err.Error()
		ctx.IndentedJSON(http.StatusBadRequest, &proto.ResponseWrapper{Code: 400, Message: "bad request", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
	}

	var order proto.Order

	if err := ctx.BindJSON(&order); err != nil {
		var errString string = err.Error()
		ctx.IndentedJSON(http.StatusBadRequest, &proto.ResponseWrapper{Code: 400, Message: "bad request", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
    }

	order.Id = id

	req := &proto.RequestWrapper{Method: getMethod(ctx), Url: getURL(ctx), RequestBody: &proto.RequestBody{Order: &order}}
	if response, err := Produce(req); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		errorHandler(ctx, 3, response)
	}
}

func PostOrder(ctx *gin.Context) {
	jwterr := proto.ValidateJWT(ctx)
	if jwterr != nil {
		errorHandler(ctx, 3, jwterr)
	}

	var order proto.Order

	if err := ctx.BindJSON(&order); err != nil {
		var errString string = err.Error()
		ctx.IndentedJSON(http.StatusBadRequest, &proto.ResponseWrapper{Code: 400, Message: "bad request", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
    }

	req := &proto.RequestWrapper{Method: getMethod(ctx), Url: getURL(ctx), RequestBody: &proto.RequestBody{Order: &order}}
	if response, err := Produce(req); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		errorHandler(ctx, 3, response)
	}
}

func PatchOrder(ctx *gin.Context) {
	jwterr := proto.ValidateJWT(ctx)
	if jwterr != nil {
		errorHandler(ctx, 3, jwterr)
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		var errString string = err.Error()
		ctx.IndentedJSON(http.StatusBadRequest, &proto.ResponseWrapper{Code: 400, Message: "bad request", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
	}

	var order proto.Order

	if err := ctx.BindJSON(&order); err != nil {
		var errString string = err.Error()
		ctx.IndentedJSON(http.StatusBadRequest, &proto.ResponseWrapper{Code: 400, Message: "bad request", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
    }

	order.Id = id

	req := &proto.RequestWrapper{Method: getMethod(ctx), Url: getURL(ctx), RequestBody: &proto.RequestBody{Order: &order}}
	if response, err := Produce(req); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		errorHandler(ctx, 3, response)
	}
}

func DeleteOrder(ctx *gin.Context) {
	jwterr := proto.ValidateJWT(ctx)
	if jwterr != nil {
		errorHandler(ctx, 3, jwterr)
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		var errString string = err.Error()
		ctx.IndentedJSON(http.StatusBadRequest, &proto.ResponseWrapper{Code: 400, Message: "bad request", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
	}

	var order proto.Order

	if err := ctx.BindJSON(&order); err != nil {
		var errString string = err.Error()
		ctx.IndentedJSON(http.StatusBadRequest, &proto.ResponseWrapper{Code: 400, Message: "bad request", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
    }

	order.Id = id

	req := &proto.RequestWrapper{Method: getMethod(ctx), Url: getURL(ctx), RequestBody: &proto.RequestBody{Order: &order}}
	if response, err := Produce(req); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		errorHandler(ctx, 3, response)
	}
}

// ERROR HANDLER

// type : 1 for user, 2 for product, 3 for order
func errorHandler(ctx *gin.Context, errType int, response *proto.ResponseWrapper) {
	switch errType{
	case 1:
		if response.GetCode() == 403 {
			ctx.IndentedJSON(http.StatusForbidden, response)
			return
		} else if response.GetCode() == 404 {
			ctx.IndentedJSON(http.StatusNotFound, response)
			return
		} else if response.GetCode() == 422 {
			ctx.IndentedJSON(http.StatusUnprocessableEntity, response)
			return
		} else if response.GetCode() == 409 {
			ctx.IndentedJSON(http.StatusConflict, response)
			return
		} else if response.GetCode() == 401 {
			ctx.IndentedJSON(http.StatusUnauthorized, response)
			return
		} else if response.GetCode() == 200 {
			ctx.IndentedJSON(http.StatusOK, response)
			return
		}
		ctx.IndentedJSON(http.StatusInternalServerError, response)
		return


	case 2:
		if response.GetCode() == 403 {
			ctx.IndentedJSON(http.StatusForbidden, response)
			return
		} else if response.GetCode() == 404 {
			ctx.IndentedJSON(http.StatusNotFound, response)
			return
		} else if response.GetCode() == 422 {
			ctx.IndentedJSON(http.StatusUnprocessableEntity, response)
			return
		} else if response.GetCode() == 409 {
			ctx.IndentedJSON(http.StatusConflict, response)
			return
		} else if response.GetCode() == 401 {
			ctx.IndentedJSON(http.StatusUnauthorized, response)
			return
		}
		ctx.IndentedJSON(http.StatusInternalServerError, response)
		return

	case 3:
		if response.GetCode() == 403 {
			ctx.IndentedJSON(http.StatusForbidden, response)
			return
		} else if response.GetCode() == 404 {
			ctx.IndentedJSON(http.StatusNotFound, response)
			return
		} else if response.GetCode() == 422 {
			ctx.IndentedJSON(http.StatusUnprocessableEntity, response)
			return
		} else if response.GetCode() == 409 {
			ctx.IndentedJSON(http.StatusConflict, response)
			return
		} else if response.GetCode() == 401 {
			ctx.IndentedJSON(http.StatusUnauthorized, response)
			return
		}
		ctx.IndentedJSON(http.StatusInternalServerError, response)
		return
	}
}

func getMethod(ctx *gin.Context) string {
	return ctx.Request.Method
}

func getURL(ctx *gin.Context) string {
	return ctx.FullPath()
}