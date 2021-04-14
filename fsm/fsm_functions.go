package fsm

import (
	"fmt"
	"time"

	"github.com/TTK4145-Students-2021/project-group_16/config"
	ec "github.com/TTK4145-Students-2021/project-group_16/elevatorController"
	"github.com/TTK4145-Students-2021/project-group_16/elevio"
)

func SetAllLocalLights(e *config.Elevator) {
	//Only sets CabRequests
	for floor := 0; floor < config.TotalFloors; floor++ {
		if e.CabRequests[floor] {
			elevio.SetButtonLamp(elevio.ButtonType(elevio.BTCab), floor, true)
		} else {
			elevio.SetButtonLamp(elevio.ButtonType(elevio.BTCab), floor, false)
		}
	}
}

func ArrivedAtFloor(e *config.Elevator, currentFloor int, doorTimer *time.Timer) {
	e.Floor = currentFloor
	elevio.SetFloorIndicator(e.Floor)
	switch e.Behaviour {
	case config.Moving:
		if ec.ShouldStop(*e) {
			elevio.SetMotorDirection(elevio.MDStop)
			elevio.SetDoorOpenLamp(true)
			ec.ClearAtCurrentFloor(e)
			SetAllLocalLights(e)
			e.Behaviour = config.DoorOpen
			doorTimer.Reset(3 * time.Second)
		}
		break
	default:
		break
	}
}

func NewOrderHandler(e *config.Elevator, newOrder elevio.ButtonEvent, channels config.FsmChannels, doorTimer *time.Timer) {
	elevio.SetButtonLamp(newOrder.Button, newOrder.Floor, true)
	switch e.Behaviour {
	case config.DoorOpen:
		if e.Floor == newOrder.Floor {
			//DoorTimer.Reset(config.DoorOpenTime)
			doorTimer.Reset(3 * time.Second)
			if newOrder.Button == elevio.BTCab {
				e.CabRequests[newOrder.Floor] = false
			} else {
				e.HallRequests[newOrder.Floor][newOrder.Button] = false
			}
		} else {
			if newOrder.Button == elevio.BTCab {
				e.CabRequests[newOrder.Floor] = true
			} else {
				e.HallRequests[newOrder.Floor][newOrder.Button] = true
			}
		}
		break
	case config.Moving:
		if newOrder.Button == elevio.BTCab {
			e.CabRequests[newOrder.Floor] = true
		} else {
			e.HallRequests[newOrder.Floor][newOrder.Button] = true
		}
		break
	case config.Idle:
		if e.Floor == newOrder.Floor && (elevio.GetFloor() != -1) {
			elevio.SetDoorOpenLamp(true)
			e.Behaviour = config.DoorOpen
			doorTimer.Reset(3 * time.Second)
		} else if e.Floor == newOrder.Floor && (elevio.GetFloor() == -1) {
			fmt.Printf("%v", "Rar Case!!!!")
			if newOrder.Button == elevio.BTCab {
				e.CabRequests[newOrder.Floor] = true
			} else {
				e.HallRequests[newOrder.Floor][newOrder.Button] = true
			}
			if e.Direction == elevio.MDUp {
				elevio.SetMotorDirection(elevio.MDDown)
			} else if e.Direction == elevio.MDDown {
				elevio.SetMotorDirection(elevio.MDUp)
			}
			e.Behaviour = config.Moving
		} else if e.Floor != newOrder.Floor && (elevio.GetFloor() == -1) {
			if newOrder.Button == elevio.BTCab {
				e.CabRequests[newOrder.Floor] = true
			} else {
				e.HallRequests[newOrder.Floor][newOrder.Button] = true
			}
			elevio.SetMotorDirection(ec.ChooseDirection(*e))
			e.Behaviour = config.Moving
		} else {
			if newOrder.Button == elevio.BTCab {
				e.CabRequests[newOrder.Floor] = true
			} else {
				e.HallRequests[newOrder.Floor][newOrder.Button] = true
			}
			e.Direction = ec.ChooseDirection(*e)
			elevio.SetMotorDirection(e.Direction)
			e.Behaviour = config.Moving

		}
		break
	}
	SetAllLocalLights(e)
}

func OnDoorTimeout(e *config.Elevator, channels config.FsmChannels, DoorTimer *time.Timer) {
	switch e.Behaviour {
	case config.DoorOpen:
		if !(e.Obstruction) {
			e.Direction = ec.ChooseDirection(*e)
			elevio.SetDoorOpenLamp(false)
			elevio.SetMotorDirection(e.Direction)
			if e.Direction == elevio.MDStop {
				e.Behaviour = config.Idle
			} else {
				e.Behaviour = config.Moving
			}
		}
		break
	default:
		break
	}
}

func ObstructionHandler(e *config.Elevator, obstruction bool, DoorTimer *time.Timer) {
	e.Obstruction = obstruction
	switch e.Behaviour {
	case config.DoorOpen:
		if !(e.Obstruction) {
			DoorTimer.Reset(3 * time.Second)
		}
		break
	default:
		break
	}
}

func StopHandler(e *config.Elevator, channels config.FsmChannels, stop bool, DoorTimer *time.Timer) {
	switch e.Behaviour {
	case config.Stop:
		elevio.SetStopLamp(false)
		if elevio.GetFloor() != -1 {
			DoorTimer.Reset(3 * time.Second)
			e.Behaviour = config.DoorOpen
		} else {
			e.Behaviour = config.Idle
		}
		break

	default:
		elevio.SetStopLamp(true)
		e.Behaviour = config.Stop
		//clearing all orders
		for f := 0; f < config.TotalFloors; f++ {
			for btn := 0; btn < config.TotalHallButtons; btn++ {
				e.HallRequests[f][btn] = false
			}
			e.CabRequests[f] = false
		}
		elevio.SetMotorDirection(elevio.MDStop)
		SetAllLocalLights(e)
		if elevio.GetFloor() != -1 {
			elevio.SetDoorOpenLamp(true)
		}
		break
	}
	SetAllLocalLights(e)
}
