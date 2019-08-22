package handlers

import (
	"net/http"

	"github.com/byuoitav/code-service/codemap"
	"github.com/labstack/echo"
)

//GetPresetHandler The endpoint to get the preset from the map
func GetPresetHandler(context echo.Context) error {
	controlKey := context.Param("controlKey")
	preset := codemap.GetPresetFromMap(controlKey)
	return context.JSON(http.StatusOK, preset)
}
