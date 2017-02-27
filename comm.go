package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
)

const MsgTypeGetConfig string = "GET_CONFIG"
const MsgTypeConfig string = "CONFIG"

type Message struct {
	Type    string `json:"type"`
	Content json.RawMessage
}

type MessageConfig struct {
	Type   string  `json:"msg"`
	Config *Config `json:"config"`
}

type ConfigRequest struct {
}

/*type SetConfig struct {
	station string
  config
}*/

func ParseMessage(msg_str []byte) (*Message, error) {

	result := &Message{}

	err := json.Unmarshal(msg_str, result)
	if err != nil {
		log.Fatalln("Can not unmarshal JSON", err)
		return nil, err
	}

	if result.Type == MsgTypeGetConfig {
		//fmt.Println("Get Config")
	} else {
		msg := fmt.Sprintf("Unknown message type '%s'", result.Type)
		log.Fatal(msg)
		return nil, err
	}

	return result, nil

}

// func ProcessMessage(reader io.Reader, writer io.Writer, cfg *Config) {
// 	io.Copy(writer, reader)
// }

func ProcessMessage(reader io.Reader, writer io.Writer, cfg *Config) {
	/*
		n, err := reader.Read(reader)

		log.Printf("msg: %s", b)

		if err != nil {
			log.Fatal("Unknown message received\n ", err)
			return
		}*/

	fmt.Printf("!\n")

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	b := buf.Bytes()

	fmt.Printf("!!=%s\n", string(b))

	//s := buf.String()

	//str_msg := string(b)

	fmt.Printf("Parsing stuff\n")

	msg_obj, err := ParseMessage(b)

	fmt.Printf(" >>> %s", string(b))

	if err != nil {
		log.Fatal("Can not parse messaeg\n ", err)
		return
	}

	// // InitWebService()
	//
	// err := json.Unmarshal(j, &colors)
	// if err != nil {
	// 	log.Fatalln("error:", err)
	// }

	//ws.Read()
	log.Printf("Received '%s'", msg_obj.Type)

	//fmt.Printf("got request... %s \n", msg_obj)

	if msg_obj.Type == MsgTypeGetConfig {
		response := MessageConfig{}
		response.Type = MsgTypeConfig

		response.Config = cfg

		data, err := json.Marshal(response)
		if err != nil {
			log.Fatalf("error while serializing: %v", err)
			//return err
		}

		data = []byte("HHHHHHHHH\n")
		log.Printf("Sending out config response, %s", data)

		sent, err := writer.Write(data)
		//writer.

		if sent < len(data) {
			fmt.Printf("Sent less")
		}

		if err != nil {
			log.Fatalf("Can not send data: %v", err)
		}

		//writer.Write(([]byte)("XXX"))

	} else {
		//fmt.Printf("Msg: %s \n", )
		log.Fatalf("Unhandled message type '%s'", msg_obj.Type)
		//io.Copy(writer, reader)
	}

}
