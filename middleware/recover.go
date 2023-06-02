package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/social-todo-list/common"
)

func Recover() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				ctx.Header("Content-Type", "application/json")

				if appErr, ok := err.(*common.AppError); ok {
					ctx.AbortWithStatusJSON(appErr.StatusCode, appErr)
					// panic(err)
					return
				}

				appErr := common.ErrInternal(err.(error))
				ctx.AbortWithStatusJSON(appErr.StatusCode, appErr)
				panic(err)
			}
		}()

		ctx.Next()
	}
}
