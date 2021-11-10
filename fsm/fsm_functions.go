package fsm

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/TTK4145-Students-2021/project-group_16/config"
	elevControl "github.com/TTK4145-Students-2021/project-group_16/elevatorController"
	"github.com/TTK4145-Students-2021/project-group_16/elevio"
)

func GetPortID() string {

	elevIDPtr := flag.Int("id", 0, "HeisID")
	flag.Parse()

	if *elevIDPtr < 15000 {
		fmt.Println("Missing or wrongly typed portID")
		os.Exit(1)
	}

	elevID := strconv.Itoa(*elevIDPtr)
	return elevID
}

func setLocalCabLights(elev *config.Elevator) {
	for floor := 0; floor < config.TotalFloors; floor++ {
		if elev.CabRequests[floor] {
			elevio.SetButtonLamp(elevio.ButtonType(elevio.BTCab), floor, true)
		} else {
			elevio.SetButtonLamp(elevio.ButtonType(elevio.BTCab), floor, false)
		}
	}
}

func arrivedAtFloor(elev *config.Elevator, currentFloor int, doorTimer *time.Timer, motorErrorTimer *time.Timer) {
	motorErrorTimer.Reset(timeUntilMotorError)
	elev.Floor = currentFloor
	elevio.SetFloorIndicator(elev.Floor)
	switch elev.Behaviour {
	case config.Moving:
		if elevControl.ShouldStop(*elev) {
			elevio.SetMotorDirection(elevio.MDStop)
			elevio.SetDoorOpenLamp(true)
			elevControl.ClearAtCurrentFloor(elev)
			setLocalCabLights(elev)
			elev.Behaviour = config.DoorOpen
			doorTimer.Reset(doorOpenDuration)
		}
	default:
		break
	}
}

func newRequestHandler(elev *config.Elevator, newRequest elevio.ButtonEvent, doorTimer *time.Timer, motorErrorTimer *time.Timer) {
	switch elev.Behaviour {
	case config.DoorOpen:
		if elev.Floor == newRequest.Floor {
			doorTimer.Reset(doorOpenDuration)
			if newRequest.Button == elevio.BTCab {
				elev.CabRequests[newRequest.Floor] = false
			} else {
				elev.LocalHallRequests[newRequest.Floor][newRequest.Button] = false
			}
		} else {
			if newRequest.Button == elevio.BTCab {
				elev.CabRequests[newRequest.Floor] = true
			} else {
				elev.LocalHallRequests[newRequest.Floor][newRequest.Button] = true
			}
		}
	case config.Moving:
		if newRequest.Button == elevio.BTCab {
			elev.CabRequests[newRequest.Floor] = true
		} else {
			elev.LocalHallRequests[newRequest.Floor][newRequest.Button] = true
		}
	case config.Idle:
		if elev.Floor == newRequest.Floor {
			elevio.SetDoorOpenLamp(true)
			elev.Behaviour = config.DoorOpen
			doorTimer.Reset(doorOpenDuration)
		} else {
			if newRequest.Button == elevio.BTCab {
				elev.CabRequests[newRequest.Floor] = true
			} else {
				elev.LocalHallRequests[newRequest.Floor][newRequest.Button] = true
			}
			elev.Direction = elevControl.ChooseDirection(*elev)
			elevio.SetMotorDirection(elev.Direction)
			elev.Behaviour = config.Moving
			motorErrorTimer.Reset(timeUntilMotorError)

		}
	}
	setLocalCabLights(elev)
}

func onDoorTimeout(elev *config.Elevator, DoorTimer *time.Timer, motorErrorTimer *time.Timer, transmitEnable chan<- bool) {
	switch elev.Behaviour {
	case config.DoorOpen:
		if !(elev.Obstruction) {
			elev.Direction = elevControl.ChooseDirection(*elev)
			elevio.SetDoorOpenLamp(false)
			elevio.SetMotorDirection(elev.Direction)
			elev.Online = true
			transmitEnable <- true
			if elev.Direction == elevio.MDStop {
				elev.Behaviour = config.Idle
			} else {
				elev.Behaviour = config.Moving
				motorErrorTimer.Reset(timeUntilMotorError)
			}
		}
	default:
		break
	}
}

func obstructionHandler(elev *config.Elevator, obstructionHigh bool, DoorTimer *time.Timer) {
	elev.Obstruction = obstructionHigh
	switch elev.Behaviour {
	case config.DoorOpen:
		if elev.Obstruction {
			DoorTimer.Reset(doorOpenDuration)
		} else {
			DoorTimer.Reset(doorOpenDuration)
		}
	default:
		break
	}
}

func motorErrorHandler(elev *config.Elevator, transmitEnable chan<- bool, motorErrorTimer *time.Timer) {
	if elevControl.HasRequestsAbove(*elev) || elevControl.HasRequestsBelow(*elev) {
		transmitEnable <- false
		elev.Online = false
		fmt.Printf("Motor error detected - going offline\n")
	} else {
		motorErrorTimer.Reset(timeUntilMotorError)
	}
}

func initializeElevator(elevID string, floorChannel <-chan int, doorTimer *time.Timer, motorErrorTimer *time.Timer) config.Elevator {
	elevio.Init("localhost:"+elevID, config.TotalFloors)

	stats := config.Stats{
		CRV: config.CVAll,
	}
	elev := config.Elevator{
		Behaviour:   config.Idle,
		Floor:       -1,
		Direction:   elevio.MDStop,
		ID:          elevID,
		Obstruction: false,
		Elevstats:   stats,
		Online:      true,
	}

	loadRequestsFromBackup(&elev)

	select {
	case currentFloor := <-floorChannel:
		fmt.Printf("currentFloor := <-fsmChannels.Floor: \n")
		fmt.Printf("This is my currentfloor: ")
		fmt.Printf("%v", currentFloor)
		arrivedAtFloor(&elev, currentFloor, doorTimer, motorErrorTimer)
		break
	default:
		elev.Behaviour = config.Moving
		motorErrorTimer.Reset(timeUntilMotorError)
		elevio.SetMotorDirection(elevio.MDDown)
	}
	return elev
}

func saveRequestsToBackup(elev config.Elevator) {
	backupRequests := BackupRequests{
		BackupLocalHallRequests: elev.LocalHallRequests,
		BackupCabRequests:       elev.CabRequests,
	}
	backupFile, _ := json.MarshalIndent(backupRequests, "", " ")
	_ = ioutil.WriteFile("backupRequests/"+elev.ID+".json", backupFile, 0644)
}

func loadRequestsFromBackup(elev *config.Elevator) {
	backupFile, err := ioutil.ReadFile("backupRequests/" + elev.ID + ".json")
	if err != nil {
		fmt.Printf("\n No backup found \n")
	} else {
		var requestsFromBackup BackupRequests
		_ = json.Unmarshal([]byte(backupFile), &requestsFromBackup)
		elev.CabRequests = requestsFromBackup.BackupCabRequests
		elev.LocalHallRequests = requestsFromBackup.BackupLocalHallRequests
		fmt.Printf("\n Requests loaded from backup \n")
	}
}

func setAllHallLights(hallRequests [config.TotalFloors][config.TotalHallButtons]bool) {
	for floor := 0; floor < config.TotalFloors; floor++ {
		for btn := 0; btn < config.TotalHallButtons; btn++ {
			if hallRequests[floor][btn] {
				elevio.SetButtonLamp(elevio.ButtonType(btn), floor, true)
			} else {
				elevio.SetButtonLamp(elevio.ButtonType(btn), floor, false)
			}
		}
	}
}
