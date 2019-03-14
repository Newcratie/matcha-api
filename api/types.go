package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"time"
)

var app App

type App struct {
	Db *sqlx.DB
	R  *gin.Engine
}

type User struct {
	Id          int16     `json:"id" db:"id"`
	Username    string    `json:"username" db:"username"`
	Email       string    `json:"email" db:"email"`
	LastName    string    `json:"lastname" db:"lastname"`
	FirstName   string    `json:"firstname" db:"firstname"`
	Password    string    `json:"password" db:"password"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	RandomToken string    `json:"random_token" db:"random_token"`
	Img1        string    `json:"img1" db:"img1"`
	Img2        string    `json:"img2" db:"img2"`
	Img3        string    `json:"img3" db:"img3"`
	Img4        string    `json:"img4" db:"img4"`
	Img5        string    `json:"img5" db:"img5"`
	Biography   string    `json:"biography" db:"biography"`
	Birthday    time.Time `json:"birthday" db:"birthday"`
	Genre       string    `json:"genre" db:"genre"`
	Interest    string    `json:"interest" db:"interest"`
	City        string    `json:"city" db:"city"`
	Zip         int       `json:"zip" db:"zip"`
	Country     string    `json:"country" db:"country"`
	Latitude    float32   `json:"latitude" db:"latitude"`
	Longitude   float32   `json:"longitude" db:"longitude"`
	GeoAllowed  bool      `json:"geo_allowed" db:"geo_allowed"`
	Online      bool      `json:"online" db:"online"`
	Rating      float32   `json:"rating" db:"rating"`
	Admin       bool      `json:"admin" db:"admin"`
	Token       string    `json:"token" db:"token"`
	AccessLvl   int       `json:"access_lvl" db:"access_lvl"`
}

type registerForm struct {
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Confirm   string    `db:"confirm"`
	Lastname  string    `db:"lastname"`
	Firstname string    `db:"firstname"`
	Birthday  time.Time `db:"birthday"`
}

type ErrorField struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type validationResponse struct {
	Valid bool         `json:"valid"`
	Fail  bool         `json:"fail"`
	Errs  []ErrorField `json:"errs"`
}

const vUsers = `(:username, :email, :lastname, :firstname, :password, :random_token, :img1, :img2, :img3, :img4, :img5, :biography, :birthday, :genre, :interest, :city, :zip, :country, :latitude, :longitude, :geo_allowed, :online, :rating, :admin)`
