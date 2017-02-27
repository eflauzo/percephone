// type TaskExecutionRecord struct {
//
// }

package main

import (
	"fmt"
	"log"
	"time"
)

import rpio "github.com/stianeikeland/go-rpio"

type State int

const (
	WAITING_SPRINKLE_TIME State = iota
	READY_SPRINKLE
	SPRINKLING
	IDLE
)

type StationControl struct {
	CurrentState State
	config       *CfgStationData
	//time_left float64 // time left to sprinkle
	sprinkle_start_time float64
}

func (station *StationControl) OutputPin(value int) {
	if station.config.Control.GPIO <= 0 {
		fmt.Printf("Fake pin ouput %d = %d\n", station.config.Control.GPIO, value)
		return
	}

	pin := rpio.Pin(station.config.Control.GPIO)

	if err := rpio.Open(); err != nil {
		panic(err)
	}

	pin.Output()
	pin.Write(rpio.State(value))

	defer rpio.Close()
}

func (station *StationControl) StartSprinkle() {
	station.OutputPin(1)
}

func (station *StationControl) StopSprinkle() {
	station.OutputPin(0)
}

func NewStationControl(cfg *CfgStationData) *StationControl {
	result := StationControl{}
	result.config = cfg
	if result.config.Enabled {
		result.CurrentState = WAITING_SPRINKLE_TIME
	} else {
		result.CurrentState = IDLE
	}
	return &result
}

type StationMaster struct {
	Stations map[string]*StationControl
}

func (master *StationMaster) ControlCycle(currenttime float64) {

	// check active stations

	active_stations := 0

	for k, v := range master.Stations {
		if v.CurrentState == SPRINKLING {
			time_left := float64(v.config.Time.Duration) - (currenttime - v.sprinkle_start_time)
			fmt.Printf("Time left: %2.2f\n", time_left)
			if time_left <= 0.0 {
				v.StartSprinkle()
				v.CurrentState = WAITING_SPRINKLE_TIME
				v.sprinkle_start_time = 0.0
				log.Println(fmt.Sprintf("Station '%s' entering 'Wait for Sprinkle Mode'\n", k))
			} else {
				active_stations += 1
			}
		}
	}

	// setting mode
	for k, v := range master.Stations {
		//fmt.Printf("Checking %s\n", k)
		if v.CurrentState == WAITING_SPRINKLE_TIME {
			//fmt.Print("   YER\n")
			// somebody waiting to sprinkle
			//v.StSprinkle()
			//v.CurrentState = SPRINKLING
			//v.sprinkle_start_time = time

			current_time_clock := time.Unix(int64(currenttime), 0)

			hour, minute, _ := current_time_clock.Clock()

			//if
			//fmt.Println("AAA")
			weekday_match := false
			for _, programmed_weekday := range v.config.Days {
				if programmed_weekday == current_time_clock.Weekday().String() {
					weekday_match = true
					//fmt.Println("1!!")
					break
				}
			}

			//fmt.Print("Hour ", hour, "   Min ", minute)

			if weekday_match && (hour == v.config.Time.Hour) && (minute == v.config.Time.Minute) {
				v.CurrentState = READY_SPRINKLE
				log.Println(fmt.Sprintf("Station '%s' entering 'Ready For Sprinkle Mode'\n", k))
				//break
			}

		}
	}

	//fmt.Printf("active_stations %d\n", active_stations)

	if active_stations == 0 {
		// can start
		for k, v := range master.Stations {
			if v.CurrentState == READY_SPRINKLE {
				// somebody waiting to sprinkle
				v.StartSprinkle()
				v.CurrentState = SPRINKLING
				v.sprinkle_start_time = currenttime
				log.Println(fmt.Sprintf("Station '%s' entering 'Sprinkling Mode'\n", k))
				break
			}
		}

	}

}

/*

// this is main worker that overlooks
type StationMaster struct {





}

*/

/*
func (*StationControl) start(){

}

func (*StationControl) stop(){

}
*/

//
// func Run() {
//
//
//
//
//
// }
