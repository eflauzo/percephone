package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

func RelayProcessRequest(ws *websocket.Conn) {
	log.Printf(">>(in)")
	ProcessMessage(ws, ws, cfg)
}

func get_config(args []interface{}, kwargs map[string]interface{}) (result *turnpike.CallResult) {

	fmt.Printf("44444444444555555")
	/*
		duration, ok := args[0].(float64)
		if !ok {
			return &turnpike.CallResult{Err: turnpike.URI("get_config.invalid-argument")}
		}
	*/

	/*
		go func() {
			time.Sleep(time.Duration(duration) * time.Second)
			client.Publish("alarm.ring", nil, nil, nil)
		}()
	*/

	data, err := json.Marshal(cfg)
	if err != nil {
		log.Fatalf("error while serializing: %v", err)
		//return err
	}

	response := string(data)

	return &turnpike.CallResult{Args: []interface{}{response}}
}

func main() {

	cfg = LoadFromFile("config.yaml")

	// InitWebService()

	// for {
	//
	// }
	//s := rpc.NewServer()
	//s.RegisterCodec(json.NewCodec(), "application/json")
	//s.RegisterService(new(HelloService), "")
	//
	// log.Fatal(
	// 	http.ListenAndServe(":8080", http.Handle("/rpc", s))
	// )

	//http.Handle("/rpc", s)

	turnpike.Debug()
	s := turnpike.NewBasicWebsocketServer("realm1")
	// server := &http.Server{
	// 		Handler: s,
	// 		Addr:    ":8000",
	// }
	// log.Println("turnpike server starting on port 8000")
	// log.Fatal(server.ListenAndServe())
	//

	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)

	http.Handle("/control", websocket.Handler(RelayProcessRequest))
	http.Handle("/centrum", s)

	//http.Handle("/", http.FileServer(http.Dir(".")))

	log.Printf("Listenning for requests")

	client, _ := s.GetLocalClient("realm1", nil)

	fmt.Printf("xxx")

	if err := client.BasicRegister("get_config", get_config); err != nil {
		panic(err)
	}

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}

	//r := mux.NewRouter()
	//r.Handle("/rpc", s)

	//http.ListenAndServe(":1234", r)

}
