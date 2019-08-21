package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

//GetPreset The endpoint to get the preset from the map
func GetPreset(context echo.Context, m map[string]string) error {
	code := context.Param("code")
	preset := m[code]

	return context.JSON(http.StatusOK, preset)
}
