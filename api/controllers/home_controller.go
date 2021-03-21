package controllers

import (
	"net/http"
	"log"
	"fmt"
	"github.com/sony/sonyflake"
	"github.com/zyahrial/blantik-be/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {

	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := flake.NextID()
	if err != nil {
		log.Fatalf("flake.NextID() failed with %s\n", err)
	}
	responses.JSON(w, http.StatusOK, id)
	fmt.Println("Your UUID is: %s", id)

}