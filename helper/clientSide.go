package helper

import (
	"net/http"

	"github.com/RakaiSeto/simple-app-may/service"
	"github.com/gin-gonic/gin"
)

// ERROR HANDLER
// type : 1 for user, 2 for product, 3 for order
func ErrorHandler(ctx *gin.Context, errType int, response *service.ResponseWrapper) {
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

func GetMethod(ctx *gin.Context) string {
	return ctx.Request.Method
}

func GetURL(ctx *gin.Context) string {
	return ctx.FullPath()
}