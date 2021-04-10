package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"errors"

	// "log"
	// "os"
	"fmt"
	// "encoding/json"

	"github.com/zyahrial/blantik-be/api/auth"
	"github.com/zyahrial/blantik-be/api/models"
	"github.com/zyahrial/blantik-be/api/responses"
	"github.com/zyahrial/blantik-be/api/utils/formaterror"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	token, err := server.SignIn(user.Email, user.Password)
	
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}

	var email = user.Email
	
	Me, err := user.FindMe(server.DB, email)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	// var me = Me{}
	fmt.Printf("%v\n", Me)

	type New_Token struct {
		Token		string		`json:"token"`
		// Product        string		`gorm:"size:255;not null" json:"product"`
	  }

	var my_token = New_Token{token}

	type Status struct {
		Success		bool		`json:"success"`
		New_Token   New_Token		`json:"token"`
		User		models.User		`json:"user"`
		// Product        string		`gorm:"size:255;not null" json:"product"`
	  }
	  
	  var p = Status{true,my_token,user}


	//   status := Status{}

	fmt.Printf("%v\n", p)

	// return status, 
	responses.JSON(w, http.StatusOK, p)
	// responses.JSON(w, http.StatusOK, my_token)
}


func (server *Server) SignIn(email, password string) (string, error) {

	var err error

	user := models.User{}

	err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}


	// type Success struct {
	// 	Status   string    `gorm:"size:255;not null;" json:"status"`
	// 	Author    User      `json:"author"`
	// 	AuthorID  uint32    `gorm:"not null" json:"author_id"`
	// 	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	// 	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	// }

	return auth.CreateToken(user.ID)

	// responses.JSON(w, http.StatusOK, p)
}

func (server *Server) Logout(w http.ResponseWriter, r *http.Request) {

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	type Status struct {
		Success		 bool		`json:"success"`
		Message        string		`json:"token"`
		// Product        string		`gorm:"size:255;not null" json:"product"`
	  }

	var message = "Logout Successfully"

	var p = Status{true,message}
	
	fmt.Printf("%v\n", uid)
	fmt.Printf("%v\n", "has logout")


	// var name = userGotten.Nickname

	// logout, err := auth.RefreshToken(user.ID)
	responses.JSON(w, http.StatusOK, p)

}