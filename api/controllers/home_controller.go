package controllers

import (
	"net/http"
	"errors"
	
	"github.com/zyahrial/blantik-be/api/responses"
	"github.com/zyahrial/blantik-be/api/auth"
	"github.com/zyahrial/blantik-be/api/models"

)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	//CHeck if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	user := models.User{}
	userGotten, err := user.FindUserByID(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	var name = userGotten.Nickname

	responses.JSON(w, http.StatusOK, "Hi")
	responses.JSON(w, http.StatusOK, name)

}