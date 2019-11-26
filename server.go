package main

import (
	"net/http"

	"github.com/byuoitav/code-service/handlers"
	"github.com/labstack/echo"
)

func main() {
	port := ":8029"
	router := echo.New()

	// Functionality Endpoints
	router.GET("/:controlKey/getPreset", handlers.GetPresetHandler)
	router.GET("/:preset/getControlKey", handlers.GetControlKeyHandler)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)

}
