package elevatorController

import (
	config "github.com/TTK4145-Students-2021/project-group_16/config"
	elevio "github.com/TTK4145-Students-2021/project-group_16/elevio"
)

func HasRequestsAbove(elev config.Elevator) bool {
	for floor := elev.Floor + 1; floor < config.TotalFloors; floor++ {
		for btn := 0; btn < config.TotalHallButtons; btn++ {
			if elev.LocalHallRequests[floor][btn] {
				return true
			}
		}
		if elev.CabRequests[floor] {
			return true
		}
	}
	return false
}

func HasRequestsBelow(elev config.Elevator) bool {
	for floor := 0; floor < elev.Floor; floor++ {
		for btn := 0; btn < config.TotalHallButtons; btn++ {
			if elev.LocalHallRequests[floor][btn] {
				return true
			}
		}
		if elev.CabRequests[floor] {
			return true
		}
	}
	return false
}

func ChooseDirection(elev config.Elevator) elevio.MotorDirection {

	switch elev.Direction {
	case elevio.MDUp:
		if HasRequestsAbove(elev) {
			return elevio.MDUp
		} else if HasRequestsBelow(elev) {
			return elevio.MDDown
		} else {
			return elevio.MDStop
		}

	case elevio.MDDown:
		if HasRequestsBelow(elev) {
			return elevio.MDDown
		} else if HasRequestsAbove(elev) {
			return elevio.MDUp
		} else {
			return elevio.MDStop
		}
	case elevio.MDStop:
		if HasRequestsAbove(elev) {
			return elevio.MDUp
		} else if HasRequestsBelow(elev) {
			return elevio.MDDown
		} else {
			return elevio.MDStop
		}

	default:
		return elevio.MDStop
	}

}

func ShouldStop(elev config.Elevator) bool {
	if elev.Floor == 0 || elev.Floor == config.TotalFloors { return true }
	switch elev.Direction {
	case elevio.MDUp:
		return (elev.LocalHallRequests[elev.Floor][elevio.BTHallUp] || elev.CabRequests[elev.Floor] || !(HasRequestsAbove(elev)))
	case elevio.MDDown:
		return (elev.LocalHallRequests[elev.Floor][elevio.BTHallDown] || elev.CabRequests[elev.Floor] || !(HasRequestsBelow(elev)))
	case elevio.MDStop:
	}
	return true
}

func ClearAtCurrentFloor(elev *config.Elevator) {
	switch elev.Elevstats.CRV {
	case config.CVAll:
		for btn := 0; btn < config.TotalHallButtons; btn++ {
			elev.LocalHallRequests[elev.Floor][btn] = false
		}
		elev.CabRequests[elev.Floor] = false
		break
	case config.CVInDirn:
		elev.CabRequests[elev.Floor] = false
		switch elev.Direction {
		case elevio.MDUp:
			elev.LocalHallRequests[elev.Floor][elevio.BTHallUp] = false
			if !HasRequestsAbove(*elev) {
				elev.LocalHallRequests[elev.Floor][elevio.BTHallDown] = false
			}
			break

		case elevio.MDDown:
			elev.LocalHallRequests[elev.Floor][elevio.BTHallDown] = false
			if !HasRequestsBelow(*elev) {
				elev.LocalHallRequests[elev.Floor][elevio.BTHallUp] = false
			}
			break

		case elevio.MDStop:
		default:
			elev.LocalHallRequests[elev.Floor][elevio.BTHallUp] = false
			elev.LocalHallRequests[elev.Floor][elevio.BTHallDown] = false
			break
		}
		break

	default:
		break
	}
}
