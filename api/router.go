package api

func (app *App) routerAPI() {
	auth := app.R.Group("/auth")
	{
		auth.POST("/login", Login)
		auth.POST("/register", Register)
	}
	api := app.R.Group("/api")
	{
		api.GET("/next", Next)
	}
}
