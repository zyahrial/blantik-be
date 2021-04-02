package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
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

	//   var me = Me{}

	type Status struct {
		Success		bool
		Token        string		`json:"Token"`
		User		 models.User			
		// Product        string		`gorm:"size:255;not null" json:"product"`
	  }

	  var p = Status{true,token,user}


	//   status := Status{}

	fmt.Printf("%v\n", Me)

	// return status, 
	responses.JSON(w, http.StatusOK, p)
	// responses.JSON(w, http.StatusOK, Me)
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
}