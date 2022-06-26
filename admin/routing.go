package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/RakaiSeto/simple-app-may/db"
	proto "github.com/RakaiSeto/simple-app-may/service"
	"github.com/gin-gonic/gin"
)

var funcCtx = context.TODO()

func Tes(ctx *gin.Context) {
	reqBody := &proto.EmptyStruct{}
	req := &proto.RequestWrapper{Method: getMethod(ctx), Url: getURL(ctx), RequestBody: &proto.RequestBody{EmptyStruct: reqBody}}
	if response, _ := Client.Tes(ctx, req); response.ResponseBody.GetError() == "" {
		ctx.IndentedJSON(http.StatusOK, response)
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

	if response, _ := Produce(req); response.ResponseBody.GetError() != "" {
		errorHandler(ctx, 1, response)
		return
	} else {
		ctx.SetCookie("token", response.ResponseBody.ResponseStatus.GetResponse(), 30*60, "/", "localhost", false, true)
		ctx.IndentedJSON(http.StatusOK, response)
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

	if response, _ := Produce(req); response.ResponseBody.GetError() != "" {
		errorHandler(ctx, 1, response)
		return
	} else {
		ctx.SetCookie("token", "", -1, "/", "localhost", false, true)
		ctx.IndentedJSON(http.StatusOK, response)
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
	if response, err := Client.AllProduct(ctx, req); err == nil {
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
	if response, err := Client.OneProduct(ctx, req); err == nil {
		ctx.IndentedJSON(http.StatusInternalServerError, response)
	} else {	
		errorHandler(ctx, 2, response)
	}
}

func PostProduct(ctx *gin.Context) {
	jwterr := proto.ValidateJWT(ctx)
	if jwterr != nil {
		errorHandler(ctx, 2, jwterr)
	}

	var adminProduct proto.AdminProduct

	if err := ctx.BindJSON(&adminProduct); err != nil {
		var errString string = err.Error()
		ctx.IndentedJSON(http.StatusBadRequest, &proto.ResponseWrapper{Code: 400, Message: "bad request", ResponseBody: &proto.ResponseBody{Error: &errString}})
		return
    }

	req := &proto.RequestWrapper{Method: getMethod(ctx), Url: getURL(ctx), RequestBody: &proto.RequestBody{AdminProduct: &adminProduct}}
	if response, err := Produce(req); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		errorHandler(ctx, 2, response)
	}
}

func PatchProduct(ctx *gin.Context) {
	var adminProduct proto.AdminProduct

	if err := ctx.BindJSON(&adminProduct); err != nil {
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

	adminProduct.Id = int64(id)

	req := &proto.RequestWrapper{Method: getMethod(ctx), Url: getURL(ctx), RequestBody: &proto.RequestBody{AdminProduct: &adminProduct}}
	if response, err := Produce(req); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		errorHandler(ctx, 2, response)
	}
}

func DeleteProduct(ctx *gin.Context) {
	var adminProduct proto.AdminProduct

	if err := ctx.BindJSON(&adminProduct); err != nil {
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

	adminProduct.Id = int64(id)

	req := &proto.RequestWrapper{Method: getMethod(ctx), Url: getURL(ctx), RequestBody: &proto.RequestBody{AdminProduct: &adminProduct}}
	if response, err := Produce(req); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
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
	if response, err := Client.AllOrder(ctx, req); err == nil {
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
	if response, err := Client.OneOrder(ctx, req); err == nil {
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

func checkJWT(ctx *gin.Context) *proto.ResponseWrapper {
	cookie, err := ctx.Cookie("token")
	if err != nil {
		var errString string = "not login yet"
		return &proto.ResponseWrapper{Code: 401, Message: "unauthorized", ResponseBody: &proto.ResponseBody{Error: &errString}}
	}

	creden, err := proto.ParseJWT(funcCtx, cookie)
	if err != nil {
		var errString string = "not login yet"
		return &proto.ResponseWrapper{Code: 401, Message: "unauthorized", ResponseBody: &proto.ResponseBody{Error: &errString}}
	}

	fmt.Println(creden)

	if int64(creden["expired"].(int64)) < int64(time.Now().Unix()) {
		var errString string = "please re-login"
		proto.DeleteJWT(ctx, creden["user"].(string), db.Rdb)
		ctx.SetCookie("token", "", -1, "/", "localhost", false, true)
		return &proto.ResponseWrapper{Code: 401, Message: "unauthorized", ResponseBody: &proto.ResponseBody{Error: &errString}}
	}

	if !proto.IsJWTExist(funcCtx, creden["user"].(string), db.Rdb){
		var errString string = "not login yet"
		return &proto.ResponseWrapper{Code: 401, Message: "unauthorized", ResponseBody: &proto.ResponseBody{Error: &errString}}
	}

	tokenString, err := proto.CheckJWT(funcCtx, creden["user"].(string), db.Rdb)
	if err != nil {
		var errString string = err.Error()
		return &proto.ResponseWrapper{Code: 401, Message: "unauthorized", ResponseBody: &proto.ResponseBody{Error: &errString}}
	}

	if *tokenString != cookie {
		var errString string = "please re-login"
		return &proto.ResponseWrapper{Code: 401, Message: "unauthorized", ResponseBody: &proto.ResponseBody{Error: &errString}}
	}
	return nil
}