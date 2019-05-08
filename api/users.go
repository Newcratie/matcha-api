package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Newcratie/matcha-api/api/hash"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"strings"
)

type Request struct {
	context *gin.Context
	claims  jwt.MapClaims
	user    User
	body    map[string]interface{}
	id      int
}

func UserModify(c *gin.Context) {
	var req Request
	if err := req.prepareRequest(c); err != nil {
		c.JSON(201, gin.H{"err": err.Error()})
	} else {
		UpdateLastConn(req.id)
		req.body = getBodyToMap(c)
		req.user, _ = app.getUser(req.id, "")
		mod := c.Param("name")
		switch mod {
		case "position":
			req.updatePosition()
			break
		case "location":
			req.updateLocation()
			break
		case "biography":
			req.updateBio()
			break
		case "username":
			req.updateUsername()
			break
		case "usertags":
			req.userTags()
			break
		case "tag":
			req.addTag()
			break
		case "password":
			req.updatePassword()
			break
		case "firstname":
			req.updateFirstname()
			break
		case "lastname":
			req.updateLastname()
			break
		case "genre":
			req.updateGenre()
			break
		case "email":
			req.updateEmail()
			break
		case "interest":
			req.updateInterest()
			break
		default:
			c.JSON(202, gin.H{"err": "route not found"})
		}
	}
}

func (req Request) updatePosition() {

	prin("POSITION ==> ", req.body, "|")

	var lat string
	var lon string

	if req.body["type"] == "ip" {
		lat, lon, _ = getPositionFromIp(req.body["position"].(string))
		prin("LAAAAAT ==> ", lat, "LOOOON ==> ", lon, "|")
	} else if req.body["type"] == "gps" {
		prin("OOOOOOOOOOOO")
	}

	//Ip := net.ParseIP(req.body["position"].(string))
	//
	//if Ip == nil {
	//	lat = req.body["lat"].(string)
	//	long = req.body["long"].(string)
	//	prin("LAT ==> ", lat, "| LONG ==> ", long, "|")
	//} else {
	//	getPositionFromIp(Ip)
	//}
	retUser(req)
}
func (req Request) updateLocation() {
	pos := req.body["location"]
	fmt.Println(pos)
	retUser(req)
}
func (req Request) updateBio() {
	bio := req.body["biography"].(string)

	if len(bio) > 100 || len(bio) < 10 {
		err := errors.New("error : your biography must be between 10 and 100 characters")
		req.context.JSON(201, gin.H{"err": err.Error()})
	} else {
		req.user.Biography = bio
		app.updateUser(req.user)
		retUser(req)
	}
}

func (req Request) checkPassword() error {
	pass := req.body["old_password"].(string)
	truePass := hash.Decrypt(hashKey, req.user.Password)
	if pass != truePass {
		return errors.New("Wrong password")
	} else {
		return nil
	}
}

func (req Request) updateUsername() {
	newUsername := req.body["new_username"].(string)
	if err := req.checkPassword(); err != nil {
		req.context.JSON(201, gin.H{"err": err.Error()})
	} else if len(newUsername) < 6 || len(newUsername) > 20 {
		err = errors.New("error : your username must be between 6 to 20 characters")
		req.context.JSON(201, gin.H{"err": err.Error()})
	} else {
		req.user.Username = newUsername
		app.updateUser(req.user)
		retUser(req)
	}
}

func (req Request) updateGenre() {
	genre := req.body["genre"].(string)
	req.user.Genre = genre
	app.updateUser(req.user)
	retUser(req)
}

func (req Request) updateInterest() {
	interest := req.body["interest"].(string)
	req.user.Interest = interest
	app.updateUser(req.user)
	retUser(req)
}

func (req Request) addTag() {
}

func (req Request) userTags() {
	tab := req.body["tags"].([]interface{})
	var userTags []string
	for _, tag := range tab {
		userTags = append(userTags, tag.(string))
	}
	fmt.Println("1", req.user.Tags)
	req.user.Tags = userTags
	fmt.Println("2", req.user.Tags)

	app.updateUser(req.user)
	retUser(req)
}

func (req Request) updateEmail() {
	newEmail := req.body["new_email"].(string)
	if err := req.checkPassword(); err != nil {
		req.context.JSON(201, gin.H{"err": err.Error()})
	} else if !emailIsValid(newEmail) {
		req.context.JSON(201, gin.H{"err": "Email is invalid"})
	} else {
		req.user.Email = newEmail
		app.updateUser(req.user)
		retUser(req)
	}
}

func (req Request) updatePassword() {
	newPassword := req.body["new_password"].(string)
	confirmPassword := req.body["confirm"].(string)

	if err := req.checkPassword(); err != nil {
		req.context.JSON(201, gin.H{"err": err.Error()})
	} else {
		if err = verifyPassword(newPassword, confirmPassword); err != nil {
			req.context.JSON(201, gin.H{"err": err.Error()})
		} else {
			req.user.Password = hash.Encrypt(hashKey, newPassword)
			app.updateUser(req.user)
			retUser(req)
		}
	}
}

func (req Request) updateFirstname() {
	firstname := req.body["firstname"].(string)
	if len(firstname) < 2 || len(firstname) > 20 {
		req.context.JSON(201, gin.H{"err": "error : your firstname must be between 2 to 20 characters"})
	} else {
		req.user.FirstName = firstname
		app.updateUser(req.user)
		retUser(req)
	}
}

func (req Request) updateLastname() {
	lastname := req.body["lastname"].(string)
	if len(lastname) < 2 || len(lastname) > 20 {
		err := errors.New("error : your lastname must be between 2 to 20 characters")
		req.context.JSON(201, gin.H{"err": err.Error()})
	} else {
		req.user.LastName = lastname
		app.updateUser(req.user)
		retUser(req)
	}
}

var magicTable = map[string]string{
	"\xff\xd8\xff":      "jpg",
	"\x89PNG\r\n\x1a\n": "png",
	"GIF87a":            "gif",
	"GIF89a":            "gif",
}

func extFromIncipit(incipit []byte) (string, error) {
	incipitStr := []byte(incipit)
	for magic, mime := range magicTable {
		if strings.HasPrefix(string(incipitStr), magic) {
			return mime, nil
		}
	}

	return "", errors.New("Wrong file")
}
func userImageHandler(c *gin.Context) {
	mFile, _ := c.FormFile("file")                // Get Multipart Header
	file, _ := mFile.Open()                       // Create Reader
	buf := bytes.NewBuffer(nil)                   // Init buffer
	if _, err := io.Copy(buf, file); err != nil { // Read file
		fmt.Println(err)
		c.JSON(201, gin.H{"err": err.Error()})
	} else {
		name := newToken() // Generate random Name
		ext, err := extFromIncipit(buf.Bytes())
		link := imageHost + "/" + name + "." + ext
		if err != nil {
			c.JSON(201, gin.H{"err": err.Error()})
			fmt.Println(err)
		} else {
			fmt.Println("ext ========> ", ext)
			f, _ := os.Create(imageSrc + "/" + name + "." + ext) //create file
			defer f.Close()                                      //close after processing

			f.Write(buf.Bytes()) // Write buffer on the file

			claims := jwt.MapClaims{}
			valid, err := ValidateToken(c, &claims)

			if valid != true {
				c.JSON(201, gin.H{"err": err.Error()})
				fmt.Println(err)
			} else {
				var req Request
				if err := req.prepareRequest(c); err != nil {
					c.JSON(201, gin.H{"err": err.Error()})
				} else {
					fmt.Println("PARAM    ", c.Param("n"))
					switch c.Param("n") {
					case "img1":
						req.user.Img1 = link
						break
					case "img2":
						req.user.Img2 = link
						break
					case "img3":
						fmt.Println("Ok     =========================")
						req.user.Img3 = link
						break
					case "img4":
						req.user.Img4 = link
						break
					case "img5":
						req.user.Img5 = link
						break
					}
					app.updateUser(req.user)
					retUser(req)
				}
			}
		}
	}
}

func getBodyToMap(c *gin.Context) (body map[string]interface{}) {
	r, _ := c.GetRawData()
	err := json.Unmarshal(r, &body)
	if err != nil {
		panic(err)
	}
	return
}

func (req *Request) prepareRequest(c *gin.Context) error {
	req.context = c
	req.claims = jwt.MapClaims{}
	valid, err := ValidateToken(c, &req.claims)
	if valid == true {
		req.id = int(req.claims["id"].(float64))
		req.user, _ = app.getUser(req.id, "")
	} else {
		return err
	}
	return nil
}

func UserHandler(c *gin.Context) {
	var req Request
	if err := req.prepareRequest(c); err != nil {
		c.JSON(201, gin.H{"err": err.Error()})
	}
	retUser(req)
}

func retUser(req Request) {
	g, err := app.dbGetUserProfile(req.id)
	tagList := app.dbGetTagList()
	userTags := app.dbGetUserTags(req.user.Username)
	if err != nil {
		req.context.JSON(201, gin.H{"err": err.Error()})
	} else {
		UpdateLastConn(req.id)
		req.context.JSON(200, gin.H{"user": g, "tagList": tagList, "userTags": userTags})
	}
}
