package codemap

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/byuoitav/common/db"
	"github.com/byuoitav/common/log"
)

var m map[string]Preset
var reqChannel chan request
var codeReqChannel chan codeRequest
var mapChannel chan map[string]Preset

type request struct {
	code   string
	respch chan Preset
}

type codeRequest struct {
	preset Preset
	respch chan ControlKey
}

// Preset struct
type Preset struct {
	RoomID     string
	PresetName string
	Ok         bool
}

//ControlKey struct
type ControlKey struct {
	ControlKey string
	Ok         bool
}

func init() {
	reqChannel = make(chan request)
	codeReqChannel = make(chan codeRequest)
	mapChannel = make(chan map[string]Preset)
	m = generateMap()
	//send events to all of the pis
	// messenger, er := messenger.BuildMessenger("", base.Messenger, 5000)
	// if er != nil {
	// 	log.L.Fatalf("failed to build messenger: %s", er)
	// }
	// // for key, value := range m {
	// // 	SendEvent(key, value.RoomID, value.PresetName, *messenger)
	// // }
	go startManager()
	go refreshMap()
}

func generateMap() map[string]Preset {
	//Query the DB for all of the UIConfigs
	uiConfigs, er := db.GetDB().GetAllUIConfigs()
	if er != nil {
		log.L.Errorf("error: %s", er)
	}
	//create a map for Room/Preset
	m = make(map[string]Preset)
	for r := range uiConfigs {
		for p := range uiConfigs[r].Presets {
			code := generateCode()
			_, exists := m[code]
			for exists {
				code = generateCode()
				_, exists = m[code]
			}
			m[code] = Preset{
				RoomID:     uiConfigs[r].ID,
				PresetName: uiConfigs[r].Presets[p].Name,
			}
		}
	}
	// //print out map
	// for key, value := range m {
	// 	fmt.Println("Key:", key, "Value:", value)
	// }
	return m
}

func generateCode() string {
	min := 0
	max := 1000000
	code := strconv.Itoa(rand.Intn(max - min))
	//prepend it with zeros so every number selected is still a 6 digit number (i.e 1234 --> 001234)
	code = "000000" + code
	code = string(code[len(code)-6:])
	return code
}

// GetPresetFromMap function
func GetPresetFromMap(code string) Preset {
	req := request{
		code:   code,
		respch: make(chan Preset),
	}
	reqChannel <- req
	return <-req.respch
}

func GetControlKeyFromPreset(preset Preset) ControlKey {
	req := codeRequest{
		preset: preset,
		respch: make(chan ControlKey),
	}
	codeReqChannel <- req
	return <-req.respch
}

func startManager() {
	for {
		select {
		case req := <-reqChannel:
			returnedPreset, ok := m[req.code]
			if !ok {
				preset := Preset{
					RoomID:     "",
					PresetName: "",
					Ok:         ok,
				}
				req.respch <- preset
			} else {
				preset := Preset{
					RoomID:     returnedPreset.RoomID,
					PresetName: returnedPreset.PresetName,
					Ok:         ok,
				}
				req.respch <- preset
			}
			close(req.respch)

		case req := <-codeReqChannel:
			counter := 0
			for key, value := range m {
				if value == req.preset {
					controlKey := ControlKey{
						ControlKey: key,
						Ok:         true,
					}
					counter = 1
					req.respch <- controlKey
					close(req.respch)
				}
			}
			if counter == 0 {
				controlKey := ControlKey{
					ControlKey: "",
					Ok:         false,
				}
				req.respch <- controlKey
				close(req.respch)
			}
		case newMap := <-mapChannel:
			m = newMap
			//send events to all of the pis
			// messenger, er := messenger.BuildMessenger("ITB-1010-CP1:7100", base.Messenger, 5000)
			// if er != nil {
			// 	log.L.Fatalf("failed to build messenger: %s", er)
			// }
			// for key, value := range m {
			// 	SendEvent(key, value.RoomID, value.PresetName, *messenger)
			// 	fmt.Println("Key:", key, "Value:", value)
			// }

		}
	}
}

// //SendEvent this emits an event that tells the pis what thier code is
// func SendEvent(controlKey string, roomID string, presetName string, runner messenger.Messenger) {
// 	a := strings.Split(roomID, "-")
// 	roominfo := events.BasicRoomInfo{}
// 	if len(a) == 2 {
// 		roominfo = events.BasicRoomInfo{
// 			BuildingID: a[0],
// 			RoomID:     roomID,
// 		}
// 	}
// 	Event := events.Event{
// 		Timestamp:    time.Now(),
// 		Key:          "ControlKey",
// 		Value:        controlKey,
// 		AffectedRoom: roominfo,
// 		EventTags: []string{
// 			events.Heartbeat,
// 		},
// 		Data: presetName,
// 	}

// 	runner.SendEvent(Event)

// }

func refreshMap() {
	ticker := time.NewTicker(1 * time.Hour)
	for range ticker.C {
		newMap := generateMap()
		mapChannel <- newMap
	}
}
