package main

import (
	"fmt"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	//stations string //`yaml:"stations"`
	Name string //`yaml:"name"`
	//Stations CfgStations `yaml:"stations"`
	Stations []CfgStation

	//control
}

type CfgStation struct {
	Station CfgStationData `yaml:"station"`
}

type CfgStationData struct {
	Name    string        `yaml:"name" json:"name"`
	Days    []string      `yaml:"days,flow" json:"days"`
	Time    CfgTargetTime `yaml:"time" json:"time"`
	Enabled bool          `yaml:"enabled" json:"enabled"`
	Control CfgControl
}

type CfgTargetTime struct {
	Hour     int `yaml:"hour" json:"hour"`
	Minute   int `yaml:"minute" json:"minute"`
	Duration int `yaml:"duration" json:"duration"`
}

type CfgControl struct {
	GPIO int `yaml:"gpio" json:"gpio"`
}

func LoadFromFile(filename string) *Config {
	data, err := ioutil.ReadFile(filename)
	fmt.Println((string)(data))
	if err != nil {
		log.Fatal(fmt.Sprintf("Can not open config file '%s'", filename))
		panic(err)
	}
	result := Config{}
	err = yaml.Unmarshal([]byte(data), &result)
	if err != nil {
		log.Fatal("Can not parse config file:", err)
		panic(err)
	}
	return &result
}

func SaveToFile(cfg *Config, filename string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		log.Fatalf("error while serializing: %v", err)
		return err
	}

	err = ioutil.WriteFile(filename, data, 0644)

	if err != nil {
		log.Fatalf("error while wiring to file: %v", err)
		return err
	}

	return nil
}
