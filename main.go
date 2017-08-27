package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
	turnpike "gopkg.in/jcelliott/turnpike.v2"
)

//import "net/http"

/*
import (
	"fmt"
	"os"
	"time"

	rpio "github.com/stianeikeland/go-rpio"
)

var (
	// Use mcu pin 10, corresponds to physical pin 19 on the pi
	//pin = rpio.Pin(23)
	pin = rpio.Pin(18)
)

func main() {
	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Unmap gpio memory when done
	defer rpio.Close()

	// Set pin to output mode
	pin.Output()

	// Toggle pin 20 times
	for x := 0; x < 20; x++ {
		pin.Toggle()
		time.Sleep(time.Second / 5)
	}
}

*/

/*
type HelloArgs struct {
	Who string
}

type HelloReply struct {
	Message string
}

type HelloService struct{}

func (h *HelloService) Say(r *http.Request, args *HelloArgs, reply *HelloReply) error {
	reply.Message = "Hello, " + args.Who + "!"
	return nil
}
*/

// func InitWebService() {
//
// }

var cfg *Config
var stationMaster *StationMaster

func RelayProcessRequest(ws *websocket.Conn) {
	log.Printf(">>(in)")
	ProcessMessage(ws, ws, cfg)
}

func get_config(args []interface{}, kwargs map[string]interface{}) (result *turnpike.CallResult) {
	data, err := json.Marshal(cfg)
	if err != nil {
		log.Fatalf("error while serializing: %v", err)
	}
	response := string(data)
	return &turnpike.CallResult{Args: []interface{}{response}}
}

func manual_start(args []interface{}, kwargs map[string]interface{}) (result *turnpike.CallResult) {
	station := args[0].(string)
	log.Printf("Manual Start: %s \n", station)
	stationMaster.Stations[station].CurrentState = READY_SPRINKLE
	return &turnpike.CallResult{Args: []interface{}{true}}
}

func manual_stop(args []interface{}, kwargs map[string]interface{}) (result *turnpike.CallResult) {
	station := args[0].(string)
	log.Printf("Manual Stop: %s \n", station)

	if stationMaster.Stations[station].config.Enabled {
		stationMaster.Stations[station].CurrentState = WAITING_SPRINKLE_TIME
	} else {
		stationMaster.Stations[station].CurrentState = IDLE
	}
	stationMaster.Stations[station].StopSprinkle()
	return &turnpike.CallResult{Args: []interface{}{true}}
}

func NewStationMaster() (master *StationMaster) {
	var result = StationMaster{}
	result.Stations = make(map[string]*StationControl)
	for i := range cfg.Stations {
		result.Stations[cfg.Stations[i].Station.Name] = NewStationControl(&cfg.Stations[i].Station)
	}
	return &result
}

type StationStatusMessage struct {
	Operation string  `yaml:"operation" json:"operation"`
	TimeLeft  float64 `yaml:"time_left" json:"time_left"`
	StartTime float64 `yaml:"start_time" json:"start_time"`
}

func operate(client *turnpike.Client, master *StationMaster) {
	for {
		var timenow = (float64)(time.Now().UTC().UnixNano() / 1000000000.0)
		fmt.Printf("Control Cycle: %2.2f\n", timenow)
		master.ControlCycle(timenow)

		stat := make(map[string]StationStatusMessage)

		for k, v := range master.Stations {
			operation := "Unknown"
			switch v.CurrentState {
			case WAITING_SPRINKLE_TIME:
				operation = "Waiting"

			case READY_SPRINKLE:
				operation = "Ready to Sprinkle"

			case SPRINKLING:
				operation = "Sprinkling"

			case IDLE:
				operation = "Idle"

			default:
				operation = "Unknown"
			}

			rec := StationStatusMessage{}
			rec.Operation = operation
			rec.TimeLeft = v.NextSprinkleTime(timenow) - timenow
			rec.StartTime = v.NextSprinkleTime(timenow)
			stat[k] = rec

			//time_left := float64(v.config.Time.Duration) - (timenow - v.sprinkle_start_time)
			fmt.Printf("    %s: %s (Time to sprinkle: %2.2f, Sprinkle time: %2.2f)\n", k, operation, rec.TimeLeft, rec.StartTime)

		}

		if err := client.Publish("station_status", nil, []interface{}{stat}, nil); err != nil {
			log.Println("Error sending message:", err)
		}

		time.Sleep(1000 * time.Millisecond)
	}
}

func main() {

	cfg = LoadFromFile("config.yaml")

	turnpike.Debug()
	s := turnpike.NewBasicWebsocketServer("realm1")

	fs := http.FileServer(http.Dir("web/build/bundled"))
	http.Handle("/", fs)

	http.Handle("/control", websocket.Handler(RelayProcessRequest))
	http.Handle("/centrum", s)

	//http.Handle("/", http.FileServer(http.Dir(".")))

	log.Printf("Listenning for requests")

	client, _ := s.GetLocalClient("realm1", nil)

	if err := client.BasicRegister("get_config", get_config); err != nil {
		panic(err)
	}

	if err := client.BasicRegister("manual_start", manual_start); err != nil {
		panic(err)
	}

	if err := client.BasicRegister("manual_stop", manual_stop); err != nil {
		panic(err)
	}

	stationMaster = NewStationMaster()
	go operate(client, stationMaster)

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}

}
