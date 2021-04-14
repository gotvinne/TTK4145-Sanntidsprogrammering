package elevatorController

import (
	config "github.com/TTK4145-Students-2021/project-group_16/config"
	elevio "github.com/TTK4145-Students-2021/project-group_16/elevio"
)

func HasOrdersAbove(e config.Elevator) bool {
	for i := e.Floor + 1; i < config.TotalFloors; i++ {
		for btn := 0; btn < config.TotalHallButtons; btn++ {
			if e.HallRequests[i][btn] {
				return true
			}
		}
		if e.CabRequests[i] {
			return true
		}
	}
	return false
}

func HasOrdersBelow(e config.Elevator) bool {
	for i := 0; i < e.Floor; i++ {
		for btn := 0; btn < config.TotalHallButtons; btn++ {
			if e.HallRequests[i][btn] {
				return true
			}
		}
		if e.CabRequests[i] {
			return true
		}
	}
	return false
}

func ChooseDirection(e config.Elevator) elevio.MotorDirection {

	switch e.Direction {
	case elevio.MDUp:
		if HasOrdersAbove(e) {
			return elevio.MDUp
		} else if HasOrdersBelow(e) {
			return elevio.MDDown
		} else {
			return elevio.MDStop
		}

	case elevio.MDDown:
		if HasOrdersBelow(e) {
			return elevio.MDDown
		} else if HasOrdersAbove(e) {
			return elevio.MDUp
		} else {
			return elevio.MDStop
		}
	case elevio.MDStop:
		if HasOrdersAbove(e) {
			return elevio.MDUp
		} else if HasOrdersBelow(e) {
			return elevio.MDDown
		} else {
			return elevio.MDStop
		}

	default:
		return elevio.MDStop
	}

}

func ShouldStop(e config.Elevator) bool {
	switch e.Direction {
	case elevio.MDUp:
		return (e.HallRequests[e.Floor][elevio.BTHallUp] || e.CabRequests[e.Floor] || !(HasOrdersAbove(e)))
	case elevio.MDDown:
		return (e.HallRequests[e.Floor][elevio.BTHallDown] || e.CabRequests[e.Floor] || !(HasOrdersBelow(e)))
	case elevio.MDStop:
	}
	return true
}

func ClearAtCurrentFloor(e *config.Elevator) {
	switch e.Elevstats.CRV {
	case config.CVAll:
		for btn := 0; btn < config.TotalHallButtons; btn++ {
			e.HallRequests[e.Floor][btn] = false
		}
		e.CabRequests[e.Floor] = false
		break
	case config.CVInDirn:
		e.CabRequests[e.Floor] = false
		switch e.Direction {
		case elevio.MDUp:
			e.HallRequests[e.Floor][elevio.BTHallUp] = false
			if !HasOrdersAbove(*e) {
				e.HallRequests[e.Floor][elevio.BTHallDown] = false
			}
			break

		case elevio.MDDown:
			e.HallRequests[e.Floor][elevio.BTHallDown] = false
			if !HasOrdersBelow(*e) {
				e.HallRequests[e.Floor][elevio.BTHallUp] = false
			}
			break

		case elevio.MDStop:
		default:
			e.HallRequests[e.Floor][elevio.BTHallUp] = false
			e.HallRequests[e.Floor][elevio.BTHallDown] = false
			break
		}
		break

	default:
		break
	}
}
