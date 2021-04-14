package fsm

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/TTK4145-Students-2021/project-group_16/config"
	"github.com/TTK4145-Students-2021/project-group_16/elevio"
)

//We use captial letter for exported functions.
func SingleElevator(elevID string, fsmChannels config.FsmChannels, commChannels config.CommChannels) {
	//Initialize elevator on network.

	elevio.Init("localhost:"+elevID, config.TotalFloors)

	stats := config.Stats{
		CRV: config.CVAll,
	}

	e := config.Elevator{
		Behaviour:   config.Idle,
		Floor:       -1,
		Direction:   elevio.MDStop,
		ID:          elevID,
		Obstruction: false,
		Stop:        false,
		Elevstats:   stats,
		Online:      true,
	}

	//Timer
	DoorTimer := time.NewTimer(config.DoorOpenTime)
	DoorTimer.Stop()
	//var atAnyFloor bool

	//Backup
	file, err := ioutil.ReadFile("backup/" + e.ID + ".json")
	if err != nil {
		_ = ioutil.WriteFile("backup/"+e.ID+".json", file, 0644)
	}

	_ = json.Unmarshal([]byte(file), &e)
	//Setter default atributter unntatt ordre
	e.Behaviour = config.Idle
	e.Floor = -1
	e.Direction = elevio.MDStop
	e.ID = elevID
	e.Obstruction = false
	e.Stop = false
	e.Elevstats = stats
	e.Online = true

	//Secure the elevator to arrive at a floor if restart occurs
	select {
	case currentFloor := <-fsmChannels.Floor:
		ArrivedAtFloor(&e, currentFloor, DoorTimer)
		break
	default:
		e.Behaviour = config.Moving
		elevio.SetMotorDirection(elevio.MDDown)
	}

	for {

		select {
		case currentFloor := <-fsmChannels.Floor:
			//fmt.Printf("%v", "currentFloor \n")
			ArrivedAtFloor(&e, currentFloor, DoorTimer)
			//config.PrintElevator(e)
			file, _ := json.MarshalIndent(e, "", " ")
			_ = ioutil.WriteFile("backup/"+e.ID+".json", file, 0644)

		case newOrder := <-fsmChannels.NewOrder:
			//fmt.Printf("%v", "newOrder \n")
			NewOrderHandler(&e, newOrder, fsmChannels, DoorTimer)
			//config.PrintElevator(e)
			file, _ := json.MarshalIndent(e, "", " ")
			_ = ioutil.WriteFile("backup/"+e.ID+".json", file, 0644)

		case <-DoorTimer.C:
			//fmt.Printf("%v", "DoorTimer \n")
			OnDoorTimeout(&e, fsmChannels, DoorTimer)
			//config.PrintElevator(e)

		case obstruction := <-fsmChannels.Obstruct:
			//fmt.Printf("%v", "Obstruction \n")
			ObstructionHandler(&e, obstruction, DoorTimer)
			//config.PrintElevator(e)

		case stop := <-fsmChannels.StopBtn:
			//fmt.Printf("%v", "Stop \n")
			StopHandler(&e, fsmChannels, stop, DoorTimer)
			//config.PrintElevator(e)
		}
		//Send current state
		go func() {
			commChannels.UpdateElevator <- e
		}()
	}
}
