package controllers

import (
	"net/http"
	"errors"
	
	"github.com/zyahrial/blantik-be/api/responses"
	"github.com/zyahrial/blantik-be/api/auth"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {

	//CHeck if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	responses.JSON(w, http.StatusOK, uid)

}