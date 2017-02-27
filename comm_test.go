package main

import (
	"bytes"
	"fmt"
	"log"
	"testing"
)

func TestParseRequest(t *testing.T) {

	strin := `{"type":"GET_CONFIG"}`

	objout, err := ParseMessage(([]byte)(strin))

	if err != nil {
		log.Fatalln("Expected successful parsing of message")
		t.Fail()
	}

	if objout.Type != MsgTypeGetConfig {
		msg := fmt.Sprintf("Expected 'get_config' but got %s", objout.Type)
		log.Fatalln(msg)
		t.Fail()
	}

	// config := LoadFromFile("test_config.yaml")
	//
	// fmt.Printf("config.name %s", config.Name)
	//
	// err := SaveToFile(config, "2.yaml")
	// if err != nil {
	// 	log.Fatal("Can not save to file", err)
	// 	//t.Fail()
	// }
	//
	// log.Fatal("x2")
	// //t.Fail()
}

func TestGetConfig(t *testing.T) {
	strin := bytes.NewBufferString(`{"type":"GET_CONFIG"}`)
	result := bytes.NewBuffer(nil)
	cfg := LoadFromFile("test_config.yaml")
	ProcessMessage(strin, result, cfg)
	//fmt.Print("++++++++++++++\n")
	fmt.Println(string(result.Bytes()))
}
