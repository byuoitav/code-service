package main

import (
	"net/http"

	"github.com/byuoitav/code-service/handlers"
	"github.com/byuoitav/common"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/v2/auth"
)

func main() {

	log.SetLevel("debug")
	port := ":8029"
	router := common.NewRouter()

	// Functionality Endpoints
	write := router.Group("", auth.AuthorizeRequest("write-state", "room", auth.LookupResourceFromAddress))
	write.GET("/:controlKey/getPreset", handlers.GetPresetHandler)
	write.GET("/:preset/getControlKey", handlers.GetControlKeyHandler)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)

}
