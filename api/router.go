package api

func (app *App) routerAPI() {
	auth := app.R.Group("/auth")
	{
		auth.POST("/login", Login)
		auth.POST("/register", Register)
		auth.GET("/valid_email/:token", Token)
	}
	api := app.R.Group("/api")
	{
		api.POST("/get_people", GetPeople)
	}
}
