package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/byuoitav/code-service-microservice/handlers"
	"github.com/byuoitav/common"
	"github.com/byuoitav/common/db"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/v2/auth"
)

func main() {
	//Query the DB for all of the UIConfigs
	uiConfigs, er := db.GetDB().GetAllUIConfigs()
	if er != nil {
		log.L.Errorf("error: %s", er)
	}
	//create a map for Room/Preset
	m := make(map[string]string)
	for r := range uiConfigs {
		for p := range uiConfigs[r].Presets {
			code := generateCode()
			_, exists := m[code]
			for exists == true {
				code = generateCode()
				_, exists = m[code]
			}
			m[code] = uiConfigs[r].ID + "/" + uiConfigs[r].Presets[p].Name
		}
	}
	//print out map
	for key, value := range m {
		fmt.Println("Key:", key, "Value:", value)
	}

	log.SetLevel("debug")
	port := ":8029"
	router := common.NewRouter()

	// Functionality Endpoints
	write := router.Group("", auth.AuthorizeRequest("write-state", "room", auth.LookupResourceFromAddress))
	// write.GET("/:roomId/:presetName/code", handlers.GetPresetCode)
	// write.GET("/:roomId/code", handlers.GetRoomCodes)
	write.Get("/:code/getPreset", handlers.GetPreset(m))

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)

}

func generateCode() string {
	min := 0
	max := 10000
	code := strconv.Itoa(rand.Intn(max - min))
	if len(code) < 4 {
		if len(code) == 1 {
			code = "000" + code
		}
		if len(code) == 2 {
			code = "00" + code
		}
		if len(code) == 3 {
			code = "0" + code
		}
	}
	return code
}

func init() {
	//Query the DB for all of the UIConfigs
	uiConfigs, er := db.GetDB().GetAllUIConfigs()
	if er != nil {
		log.L.Errorf("error: %s", er)
	}
	//create a map for Room/Preset
	m := make(map[string]string)
	for r := range uiConfigs {
		for p := range uiConfigs[r].Presets {
			code := generateCode()
			_, exists := m[code]
			for exists == true {
				code = generateCode()
				_, exists = m[code]
			}
			m[code] = uiConfigs[r].ID + "/" + uiConfigs[r].Presets[p].Name
		}
	}
	//print out map
	for key, value := range m {
		fmt.Println("Key:", key, "Value:", value)
	}
}
