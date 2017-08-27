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
		//panic(err)
		log.Printf("Can not open GPIO PIN %s", err)
		return
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

func NewStationControl(station_cfg *CfgStationData) *StationControl {
	new_station_control := new(StationControl)
	new_station_control.config = station_cfg
	if new_station_control.config.Enabled {
		new_station_control.CurrentState = WAITING_SPRINKLE_TIME
	} else {
		new_station_control.CurrentState = IDLE
	}
	return new_station_control
}

type StationMaster struct {
	Stations map[string]*StationControl
}

func (station *StationControl) NextSprinkleTime(currenttime float64) float64 {
	// we iterate of days starting from this day, and check times for sprinkle
	// if sprinkle in future we return
	current_time_clock := time.Unix(int64(currenttime), 0)
	year, month, day := current_time_clock.Date()

	//wd := current_time_clock.Weekday()
	for day_offset := 0; day_offset < 8; day_offset++ {

		sprinkle_time := time.Date(
			year,
			month,
			day,
			station.config.Time.Hour,
			station.config.Time.Minute,
			0,
			0,
			current_time_clock.Location()).Add(time.Duration(day_offset*24) * time.Hour)

		weekday_match := false
		for _, programmed_weekday := range station.config.Days {
			if programmed_weekday == sprinkle_time.Weekday().String() {
				weekday_match = true
				break
			}
		}

		if weekday_match {
			unix_sprinkle_time := (float64)(sprinkle_time.Unix())
			if unix_sprinkle_time > currenttime {
				return unix_sprinkle_time
			}
		}
	}
	return -1.0
}

func (master *StationMaster) ControlCycle(currenttime float64) {

	// check active stations

	active_stations := 0

	for k, v := range master.Stations {
		if v.CurrentState == SPRINKLING {
			time_left := float64(v.config.Time.Duration) - (currenttime - v.sprinkle_start_time)
			fmt.Printf("Time left: %2.2f\n", time_left)
			if time_left <= 0.0 {
				v.StopSprinkle()
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

			current_time_clock := time.Unix(int64(currenttime), 0)

			hour, minute, _ := current_time_clock.Clock()
			weekday_match := false
			for _, programmed_weekday := range v.config.Days {
				if programmed_weekday == current_time_clock.Weekday().String() {
					weekday_match = true
					break
				}
			}

			if weekday_match && (hour == v.config.Time.Hour) && (minute == v.config.Time.Minute) {
				v.CurrentState = READY_SPRINKLE
				log.Println(fmt.Sprintf("Station '%s' entering 'Ready For Sprinkle Mode'\n", k))
				//break
			}

		}
	}

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
