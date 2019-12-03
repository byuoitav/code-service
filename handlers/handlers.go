package handlers

import (
	"net/http"
	"strings"

	"github.com/byuoitav/code-service/codemap"
	"github.com/labstack/echo"
)

//GetPresetHandler The endpoint to get the preset from the map
func GetPresetHandler(context echo.Context) error {
	controlKey := context.Param("controlKey")
	preset := codemap.GetPresetFromMap(controlKey)
	if !preset.Ok {
		return context.JSON(http.StatusNotFound, "The preset was not found for this control key")
	}
	return context.JSON(http.StatusOK, preset)
}

func GetControlKeyHandler(context echo.Context) error {
	presetParam := context.Param("preset")
	presetParts := strings.SplitN(presetParam, " ", 2)
	preset := codemap.Preset{
		RoomID:     presetParts[0],
		PresetName: presetParts[1],
	}
	controlKey := codemap.GetControlKeyFromPreset(preset)
	if !controlKey.Ok {
		return context.JSON(http.StatusNotFound, "The control key was not found for this preset")
	}
	return context.JSON(http.StatusOK, controlKey)
}

func HealthCheck(context echo.Context) error {
	return context.JSON(http.StatusOK, "Healthy!")
}
