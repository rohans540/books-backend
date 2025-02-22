package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rohans540/books-backend/controllers"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/books")
	{
		api.GET("", controllers.GetBooks)
		api.GET("/:id", controllers.GetBookByID)
		api.POST("", controllers.CreateBook)
		api.PUT("/:id", controllers.UpdateBook)
		api.DELETE("/:id", controllers.DeleteBook)
	}
}
