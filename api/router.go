package api

import (
	"github.com/gin-gonic/gin"
)

func (app *App) routerAPI() {
	auth := app.R.Group("/auth")
	{
		auth.POST("/login", Login)
		auth.POST("/register", Register)
	}
	admin := app.R.Group("/admin")
	{
		admin.POST("/users", Users)
		admin.GET("/test", func(c *gin.Context) {
			app, _ := c.MustGet("app").(App)
			c.JSON(200, app.Users)
		})
	}
	api := app.R.Group("/api")
	{
		api.GET("/next", Next)
	}
}
