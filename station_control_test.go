package main

import "testing"

func TestStationMasterControl(t *testing.T) {
	/*
		config := LoadFromFile("test_config.yaml")

		fmt.Printf("config.name %s", config.Name)

		err := SaveToFile(config, "2.yaml")
		if err != nil {
			log.Fatal("Can not save to file", err)
			t.Fail()
		}

		log.Fatal("x2")
		t.Fail()
	*/

	control := CfgControl{}
	control.GPIO = 0

	targettime := CfgTargetTime{}
	targettime.Hour = 21
	targettime.Minute = 30
	targettime.Duration = 120

	station1cfg := CfgStationData{}
	station1cfg.Name = "Station#1"
	station1cfg.Days = []string{"Sunday", "Monday"}
	station1cfg.Time = targettime
	station1cfg.Enabled = true
	station1cfg.Control = control

	station2cfg := CfgStationData{}
	station2cfg.Name = "Station#2"
	station2cfg.Days = []string{"Sunday", "Monday"}
	station2cfg.Time = targettime
	station2cfg.Enabled = true
	station2cfg.Control = control

	stationControl1 := NewStationControl(&station1cfg)
	stationControl2 := NewStationControl(&station2cfg)

	stationMaster := StationMaster{}
	stationMaster.Stations = map[string]*StationControl{
		"Station#1": stationControl1,
		"Station#2": stationControl2,
	}

	stationMaster.ControlCycle(0.0)
	if stationControl1.CurrentState != WAITING_SPRINKLE_TIME {
		t.Fatal("Expected 'wait for sprinkle time' state for station1 ")
		t.Fail()
	}

	if stationControl2.CurrentState != WAITING_SPRINKLE_TIME {
		t.Fatal("Expected 'wait for sprinkle time' state for station2")
		t.Fail()
	}

	// 11/21/2016, 9:30:00 PM (Monday)
	stationMaster.ControlCycle(1479792600.0)

	var first_sprinkling_station *StationControl = nil
	var second_sprinkling_station *StationControl = nil

	if stationControl1.CurrentState == SPRINKLING {
		first_sprinkling_station = stationControl1
		second_sprinkling_station = stationControl2
	} else if stationControl2.CurrentState == SPRINKLING {
		second_sprinkling_station = stationControl1
		first_sprinkling_station = stationControl2
	} else {
		t.Fatal("Expect one station to sprinkle")
		t.Fail()
	}

	if first_sprinkling_station.CurrentState != SPRINKLING {
		t.Fatal("Expected 'wait for sprinkle time' state for station1 ")
		t.Fail()
	}

	if second_sprinkling_station.CurrentState != READY_SPRINKLE {
		t.Fatal("Expected 'ready sprinkle' state for station2")
		t.Fail()
	}

	// +30 sec
	// expect no changes
	stationMaster.ControlCycle(1479792630.0)

	if first_sprinkling_station.CurrentState != SPRINKLING {
		t.Fatal("Expected 'wait for sprinkle time' state for station1 ")
		t.Fail()
	}

	if second_sprinkling_station.CurrentState != READY_SPRINKLE {
		t.Fatal("Expected 'ready sprinkle' state for station2")
		t.Fail()
	}

	// 11/21/2016, 9:32:00 PM (Monday)
	// T+2 min sprinkler1  should shut down
	// sprinkler2 should start

	stationMaster.ControlCycle(1479792720.0)

	if first_sprinkling_station.CurrentState != WAITING_SPRINKLE_TIME {
		t.Fatal("Expected 'wait for sprinkle time' state for station1 ")
		t.Fail()
	}

	if second_sprinkling_station.CurrentState != SPRINKLING {
		t.Fatal("Expected 'ready sprinkle' state for station2")
		t.Fail()
	}

	// 119 sec in sprinkling, still expected sprinkler1 waiting
	// and sprinkler2 sprinkling

	stationMaster.ControlCycle(1479792839.0)

	if first_sprinkling_station.CurrentState != WAITING_SPRINKLE_TIME {
		t.Fatal("Expected 'wait for sprinkle time' state for station1 ")
		t.Fail()
	}

	if second_sprinkling_station.CurrentState != SPRINKLING {
		t.Fatal("Expected 'sprinkle' state for station2")
		t.Fail()
	}

	// 120 sec in sprinkling, both stations should shutdown

	stationMaster.ControlCycle(1479792840.0)

	if first_sprinkling_station.CurrentState != WAITING_SPRINKLE_TIME {
		t.Fatal("Expected 'wait for sprinkle time' state for station1 ")
		t.Fail()
	}

	if second_sprinkling_station.CurrentState != WAITING_SPRINKLE_TIME {
		t.Fatal("Expected 'wait for sprinkle time' state for station2")
		t.Fail()
	}

	//t.Pass()
}
