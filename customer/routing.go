package main

import (
	"net/http"
	"strconv"
	proto "github.com/RakaiSeto/simple-app-may/service"
	"github.com/RakaiSeto/simple-app-may/helper"
	"github.com/gin-gonic/gin"
)


func Tes(ctx *gin.Context) {
	reqBody := &proto.EmptyStruct{}
	req := &proto.RequestWrapper{Method: helper.GetMethod(ctx), Url: helper.GetURL(ctx), RequestBody: &proto.RequestBody{EmptyStruct: reqBody}}
	if response, _ := Client.Tes(ctx, req); response.ResponseBody.GetError() == "" {
		ctx.IndentedJSON(http.StatusOK, response)
		return
	} else {
		helper.ErrorHandler(ctx, 1, response)
        return
	}
}

func Login(ctx *gin.Context) {
	var user proto.User

	if err := ctx.BindJSON(&user); err != nil {
		var errString string = err.Error()
		ctx.IndentedJSON(http.StatusBadRequest, &proto.ResponseWrapper{Code: 400, Message: "bad request", ResponseBody: &proto.ResponseBody{Error: &errString}})
	}

	req := &proto.RequestWrapper{Method: helper.GetMethod(ctx), Url: helper.GetURL(ctx), RequestBody: &proto.RequestBody{User: &user}}

	if response, _ := Client.Login(ctx, req.RequestBody); response.ResponseBody.GetError() == "" {
		ctx.SetCookie("token", response.ResponseBody.GetString_(), 30*60, "/", "localhost", false, true)
		response.ResponseBody.String_ = nil
		ctx.IndentedJSON(http.StatusOK, response)
	} else {
		helper.ErrorHandler(ctx, 1, response)
		return
	}
}

func LoginGithub(ctx *gin.Context) {
	cookie, _ := ctx.Cookie("token")
	if cookie != "" {
		var errString string = "please log out first"
		helper.ErrorHandler(ctx, 1, &proto.ResponseWrapper{Code: 400, Message: "bad request", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
	}

	reqBody := &proto.EmptyStruct{}
	req := &proto.RequestWrapper{Method: helper.GetMethod(ctx), Url: helper.GetURL(ctx), RequestBody: &proto.RequestBody{EmptyStruct: reqBody}}
	if response, _ := Client.LoginGithub(ctx, req); response.ResponseBody.GetError() == "" {
		ctx.Redirect(http.StatusTemporaryRedirect, response.ResponseBody.ResponseStatus.GetResponse())
		return
	} else {
		helper.ErrorHandler(ctx, 1, response)
        return
	}
}

func LoginGithubCallback(ctx *gin.Context) {
	state := ctx.Query("state")
	code := ctx.Query("code")
	req := &proto.RequestWrapper{Method: helper.GetMethod(ctx), Url: helper.GetURL(ctx), RequestBody: &proto.RequestBody{OauthCallback: &proto.OauthCallback{State: state, Code: code}}}
	if response, _ := Client.LoginGithubCallback(ctx, req); response.ResponseBody.GetError() == "" {
		ctx.SetCookie("token", response.ResponseBody.GetString_(), 30*60, "/", "localhost", false, true)
		response.ResponseBody.String_ = nil
		ctx.IndentedJSON(http.StatusOK, response)
		return
	} else {
		helper.ErrorHandler(ctx, 1, response)
        return
	}
}

func LoginGoogle(ctx *gin.Context) {
	cookie, _ := ctx.Cookie("token")
	if cookie != "" {
		var errString string = "please log out first"
		helper.ErrorHandler(ctx, 1, &proto.ResponseWrapper{Code: 400, Message: "bad request", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
	}

	reqBody := &proto.EmptyStruct{}
	req := &proto.RequestWrapper{Method: helper.GetMethod(ctx), Url: helper.GetURL(ctx), RequestBody: &proto.RequestBody{EmptyStruct: reqBody}}
	if response, _ := Client.LoginGoogle(ctx, req); response.ResponseBody.GetError() == "" {
		ctx.Redirect(http.StatusTemporaryRedirect, response.ResponseBody.ResponseStatus.GetResponse())
		return
	} else {
		helper.ErrorHandler(ctx, 1, response)
        return
	}
}

func LoginGoogleCallback(ctx *gin.Context) {
	state := ctx.Query("state")
	code := ctx.Query("code")
	req := &proto.RequestWrapper{Method: helper.GetMethod(ctx), Url: helper.GetURL(ctx), RequestBody: &proto.RequestBody{OauthCallback: &proto.OauthCallback{State: state, Code: code}}}
	if response, _ := Client.LoginGoogleCallback(ctx, req); response.ResponseBody.GetError() == "" {
		ctx.SetCookie("token", response.ResponseBody.GetString_(), 30*60, "/", "localhost", false, true)
		response.ResponseBody.String_ = nil
		ctx.IndentedJSON(http.StatusOK, response)
		return
	} else {
		helper.ErrorHandler(ctx, 1, response)
        return
	}
}

func LoginFacebook(ctx *gin.Context) {
	cookie, _ := ctx.Cookie("token")
	if cookie != "" {
		var errString string = "please log out first"
		helper.ErrorHandler(ctx, 1, &proto.ResponseWrapper{Code: 400, Message: "bad request", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
	}

	reqBody := &proto.EmptyStruct{}
	req := &proto.RequestWrapper{Method: helper.GetMethod(ctx), Url: helper.GetURL(ctx), RequestBody: &proto.RequestBody{EmptyStruct: reqBody}}
	if response, _ := Client.LoginFacebook(ctx, req); response.ResponseBody.GetError() == "" {
		ctx.Redirect(http.StatusTemporaryRedirect, response.ResponseBody.ResponseStatus.GetResponse())
		return
	} else {
		helper.ErrorHandler(ctx, 1, response)
        return
	}
}

func LoginFacebookCallback(ctx *gin.Context) {
	state := ctx.Query("state")
	code := ctx.Query("code")
	req := &proto.RequestWrapper{Method: helper.GetMethod(ctx), Url: helper.GetURL(ctx), RequestBody: &proto.RequestBody{OauthCallback: &proto.OauthCallback{State: state, Code: code}}}
	if response, _ := Client.LoginFacebookCallback(ctx, req); response.ResponseBody.GetError() == "" {
		ctx.SetCookie("token", response.ResponseBody.GetString_(), 30*60, "/", "localhost", false, true)
		response.ResponseBody.String_ = nil
		ctx.IndentedJSON(http.StatusOK, response)
		return
	} else {
		helper.ErrorHandler(ctx, 1, response)
        return
	}
}

func Logout(ctx *gin.Context) {
	cookie, err := ctx.Cookie("token")
	if err != nil {
		var errString string = "not logged in yet"
		helper.ErrorHandler(ctx, 1, &proto.ResponseWrapper{Code: 422, Message: "unprocessable entity", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
	}

	req := &proto.RequestBody{String_: &cookie}

	if response, _ := Client.Logout(ctx, req); response.ResponseBody.GetError() != "" {
		helper.ErrorHandler(ctx, 1, response)
		return
	} else {
		ctx.SetCookie("token", "", -1, "/", "localhost", false, true)
		ctx.IndentedJSON(http.StatusOK, response)
	}
}

func CheckError(ctx *gin.Context) {
	id := ctx.Param("id")
	req := &proto.RequestWrapper{Method: helper.GetMethod(ctx), Url: helper.GetURL(ctx), RequestBody: &proto.RequestBody{String_: &id}}
	if response, _ := Client.ReqError(ctx, req); response.ResponseBody.GetError() == "" {
		ctx.IndentedJSON(http.StatusOK, response)
		return
	} else {
		helper.ErrorHandler(ctx, 1, response)
	}
}

func AllUser(ctx *gin.Context) {
	jwterr := helper.ValidateJWT(ctx)
	if jwterr != nil {
		helper.ErrorHandler(ctx, 1, jwterr)
		return
	}

	reqBody := &proto.EmptyStruct{}
	req := &proto.RequestWrapper{Method: helper.GetMethod(ctx), Url: helper.GetURL(ctx), RequestBody: &proto.RequestBody{EmptyStruct: reqBody}}
	if response, _ := Client.AllUser(ctx, req); response.ResponseBody.GetError() == "" {
		ctx.IndentedJSON(http.StatusOK, response)
		return
	} else {	
		helper.ErrorHandler(ctx,  1, response)
		return
	}
}

func OneUser(ctx *gin.Context) {
	jwterr := helper.ValidateJWT(ctx)
	if jwterr != nil {
		helper.ErrorHandler(ctx, 1, jwterr)
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		var errString string = err.Error()
		ctx.IndentedJSON(http.StatusBadRequest, &proto.ResponseWrapper{Code: 400, Message: "bad request", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
	}

	req := &proto.RequestWrapper{Method: helper.GetMethod(ctx), Url: helper.GetURL(ctx), RequestBody: &proto.RequestBody{Id: &id}}
	if response, _ := Client.OneUser(ctx, req); response.ResponseBody.GetError() != "" {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		helper.ErrorHandler(ctx, 1, response)
	}
}

func MyUser(ctx *gin.Context) {
	jwterr := helper.ValidateJWT(ctx)
	if jwterr != nil {
		helper.ErrorHandler(ctx, 1, jwterr)
		return
	}

	token, _ := ctx.Cookie("token")
	
	req := &proto.RequestWrapper{Method: helper.GetMethod(ctx), Url: helper.GetURL(ctx), RequestBody: &proto.RequestBody{String_: &token}}
	if response, _ := Client.MyUser(ctx, req); response.ResponseBody.GetError() == "" {
		ctx.IndentedJSON(http.StatusOK, response)
		return
	} else {	
		helper.ErrorHandler(ctx,  1, response)
		return
	}
}

func PostUser(ctx *gin.Context) {
	var user proto.User

	if err := ctx.BindJSON(&user); err != nil {
		var errString string = err.Error()
		ctx.IndentedJSON(http.StatusInternalServerError, &proto.ResponseWrapper{Code: 500, Message: "unknown error", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
	}

	req := &proto.RequestBody{User: &user}
	if response, _ := Client.AddUser(ctx, req); response.ResponseBody.GetError() != "" {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		helper.ErrorHandler(ctx, 1, response)
	}
}

func PatchUser(ctx *gin.Context) {
	jwterr := helper.ValidateJWT(ctx)
	if jwterr != nil {
		helper.ErrorHandler(ctx, 1, jwterr)
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
	token, _ := ctx.Cookie("token")

	req := &proto.RequestWrapper{Method: helper.GetMethod(ctx), Url: helper.GetURL(ctx), RequestBody: &proto.RequestBody{User: &user, String_: &token}}
	if response, _ := Produce(req); response.ResponseBody.GetError() != "" {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		helper.ErrorHandler(ctx, 1, response)
	}
}

func DeleteUser(ctx *gin.Context) {
	jwterr := helper.ValidateJWT(ctx)
	if jwterr != nil {
		helper.ErrorHandler(ctx, 1, jwterr)
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

	req := &proto.RequestWrapper{Method: helper.GetMethod(ctx), Url: helper.GetURL(ctx), RequestBody: &proto.RequestBody{User: &user}}
	if response, _ := Produce(req); response.ResponseBody.GetError() != "" {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		helper.ErrorHandler(ctx, 1, response)
	}
}

func AllProduct(ctx *gin.Context) {
	jwterr := helper.ValidateJWT(ctx)
	if jwterr != nil {
		helper.ErrorHandler(ctx, 2, jwterr)
		return
	}

	reqbody := &proto.EmptyStruct{}
	req := &proto.RequestWrapper{Method: helper.GetMethod(ctx), Url: helper.GetURL(ctx), RequestBody: &proto.RequestBody{EmptyStruct: reqbody}}
	if response, _ := Client.AllProduct(ctx, req); response.ResponseBody.GetError() == "" {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		helper.ErrorHandler(ctx, 2, response)
	}
}

func OneProduct(ctx *gin.Context) {
	jwterr := helper.ValidateJWT(ctx)
	if jwterr != nil {
		helper.ErrorHandler(ctx, 2, jwterr)
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		var errString string = err.Error()
		ctx.IndentedJSON(http.StatusBadRequest, &proto.ResponseWrapper{Code: 400, Message: "bad request", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
	}

	req := &proto.RequestWrapper{Method: helper.GetMethod(ctx), Url: helper.GetURL(ctx), RequestBody: &proto.RequestBody{Id: &id}}
	if response, _ := Client.OneProduct(ctx, req); response.ResponseBody.GetError() == "" {
		ctx.IndentedJSON(http.StatusInternalServerError, response)
	} else {	
		helper.ErrorHandler(ctx, 2, response)
	}
}

func PostProduct(ctx *gin.Context) {
	jwterr := helper.ValidateJWT(ctx)
	if jwterr != nil {
		helper.ErrorHandler(ctx, 3, jwterr)
	}

	cookie, _ := ctx.Cookie("token")

	creden, err := helper.ParseJWT(ctx, cookie)
	if err != nil {
		errString := "please re-login"
		ctx.IndentedJSON(http.StatusUnauthorized, &proto.ResponseWrapper{Code: 401, Message: "unauthorized", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
	}

	if !creden["admin"].(bool){
		errString := "you're not an admin"
		ctx.IndentedJSON(http.StatusUnauthorized, &proto.ResponseWrapper{Code: 401, Message: "unauthorized", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
	}

	var product proto.Product

	if err := ctx.BindJSON(&product); err != nil {
		var errString string = err.Error()
		ctx.IndentedJSON(http.StatusBadRequest, &proto.ResponseWrapper{Code: 400, Message: "bad request", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
    }

	req := &proto.RequestWrapper{Method: helper.GetMethod(ctx), Url: helper.GetURL(ctx), RequestBody: &proto.RequestBody{Product: &product}}
	if response, err := Produce(req); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		helper.ErrorHandler(ctx, 2, response)
	}
}

func PatchProduct(ctx *gin.Context) {
	jwterr := helper.ValidateJWT(ctx)
	if jwterr != nil {
		helper.ErrorHandler(ctx, 3, jwterr)
	}

	cookie, _ := ctx.Cookie("token")

	creden, err := helper.ParseJWT(ctx, cookie)
	if err != nil {
		errString := "please re-login"
		ctx.IndentedJSON(http.StatusUnauthorized, &proto.ResponseWrapper{Code: 401, Message: "unauthorized", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
	}

	if !creden["admin"].(bool){
		errString := "you're not an admin"
		ctx.IndentedJSON(http.StatusUnauthorized, &proto.ResponseWrapper{Code: 401, Message: "unauthorized", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
	}

	var product proto.Product

	if err := ctx.BindJSON(&product); err != nil {
		var errString string = err.Error()
		ctx.IndentedJSON(http.StatusBadRequest, &proto.ResponseWrapper{Code: 400, Message: "bad request", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
    }

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		var errString string = err.Error()
		ctx.IndentedJSON(http.StatusBadRequest, &proto.ResponseWrapper{Code: 400, Message: "bad request", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
	}

	product.Id = int64(id)

	req := &proto.RequestWrapper{Method: helper.GetMethod(ctx), Url: helper.GetURL(ctx), RequestBody: &proto.RequestBody{Product: &product}}
	if response, err := Produce(req); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		helper.ErrorHandler(ctx, 2, response)
	}
}

func DeleteProduct(ctx *gin.Context) {
	jwterr := helper.ValidateJWT(ctx)
	if jwterr != nil {
		helper.ErrorHandler(ctx, 3, jwterr)
	}

	cookie, _ := ctx.Cookie("token")

	creden, err := helper.ParseJWT(ctx, cookie)
	if err != nil {
		errString := "please re-login"
		ctx.IndentedJSON(http.StatusUnauthorized, &proto.ResponseWrapper{Code: 401, Message: "unauthorized", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
	}

	if !creden["admin"].(bool){
		errString := "you're not an admin"
		ctx.IndentedJSON(http.StatusUnauthorized, &proto.ResponseWrapper{Code: 401, Message: "unauthorized", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
	}

	var product proto.Product

	if err := ctx.BindJSON(&product); err != nil {
		var errString string = err.Error()
		ctx.IndentedJSON(http.StatusBadRequest, &proto.ResponseWrapper{Code: 400, Message: "bad request", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
    }

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		var errString string = err.Error()
		ctx.IndentedJSON(http.StatusBadRequest, &proto.ResponseWrapper{Code: 400, Message: "bad request", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
	}

	product.Id = int64(id)

	req := &proto.RequestWrapper{Method: helper.GetMethod(ctx), Url: helper.GetURL(ctx), RequestBody: &proto.RequestBody{Product: &product}}
	if response, err := Produce(req); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		helper.ErrorHandler(ctx, 2, response)
	}
}

func AllOrder(ctx *gin.Context) {
	jwterr := helper.ValidateJWT(ctx)
	if jwterr != nil {
		helper.ErrorHandler(ctx, 3, jwterr)
	}

	cookie, _ := ctx.Cookie("token")

	req := &proto.RequestWrapper{Method: helper.GetMethod(ctx), Url: helper.GetURL(ctx), RequestBody: &proto.RequestBody{String_: &cookie}}
	if response, _ := Client.AllOrder(ctx, req); response.ResponseBody.GetError() == "" {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {
		helper.ErrorHandler(ctx, 3, response)
	}
}

func OneOrder(ctx *gin.Context) {
	jwterr := helper.ValidateJWT(ctx)
	if jwterr != nil {
		helper.ErrorHandler(ctx, 3, jwterr)
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		var errString string = err.Error()
		ctx.IndentedJSON(http.StatusBadRequest, &proto.ResponseWrapper{Code: 400, Message: "bad request", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
	}
	cookie, _ := ctx.Cookie("token")

	req := &proto.RequestWrapper{Method: helper.GetMethod(ctx), Url: helper.GetURL(ctx), RequestBody: &proto.RequestBody{Id: &id, String_: &cookie}}
	if response, _ := Client.OneOrder(ctx, req); response.ResponseBody.GetError() == "" {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		helper.ErrorHandler(ctx, 3, response)
	}
}

func PostOrder(ctx *gin.Context) {
	jwterr := helper.ValidateJWT(ctx)
	if jwterr != nil {
		helper.ErrorHandler(ctx, 3, jwterr)
	}

	cookie, _ := ctx.Cookie("token")

	var order proto.Order

	if err := ctx.BindJSON(&order); err != nil {
		var errString string = err.Error()
		ctx.IndentedJSON(http.StatusBadRequest, &proto.ResponseWrapper{Code: 400, Message: "bad request", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
    }

	req := &proto.RequestWrapper{Method: helper.GetMethod(ctx), Url: helper.GetURL(ctx), RequestBody: &proto.RequestBody{String_: &cookie, Order: &order}}
	if response, _ := Produce(req); response.ResponseBody.GetError() == "" {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		helper.ErrorHandler(ctx, 3, response)
	}
}

func DeleteOrder(ctx *gin.Context) {
	jwterr := helper.ValidateJWT(ctx)
	if jwterr != nil {
		helper.ErrorHandler(ctx, 3, jwterr)
	}

	cookie, _ := ctx.Cookie("token")

	var order proto.Order

	order.Id = ctx.Param("id")

	req := &proto.RequestWrapper{Method: helper.GetMethod(ctx), Url: helper.GetURL(ctx), RequestBody: &proto.RequestBody{Order: &order, String_: &cookie}}
	if response, _ := Produce(req); response.ResponseBody.GetError() == "" {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		helper.ErrorHandler(ctx, 3, response)
	}
}