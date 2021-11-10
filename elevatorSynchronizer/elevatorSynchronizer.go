package elevatorSynchronizer

import (
	"github.com/TTK4145-Students-2021/project-group_16/config"
	"github.com/TTK4145-Students-2021/project-group_16/elevio"
)

type ElevSyncChannels struct {
	IncomingRequest   chan elevio.ButtonEvent
	IncomingElevator  chan config.Elevator
	PeersLost         chan []string
	BecomeIndependant chan bool
}

type CostRequest struct {
	HallRequests    [config.TotalFloors][config.TotalHallButtons]bool `json:"hallRequests"`
	OnlineElevators map[string]config.Elevator                        `json:"states"`
}

func ElevatorSynchronizer(elevID string, elevSyncChannels ElevSyncChannels, fsmChannels config.FsmChannels) {

	elevatorMap := make(map[string]config.Elevator)

	for {
		globalHallRequests := [config.TotalFloors][config.TotalHallButtons]bool{}

		select {
		case incomingRequest := <-elevSyncChannels.IncomingRequest:

			if incomingRequest.Button == elevio.BTCab {
				fsmChannels.NewRequest <- incomingRequest
			} else {
				globalHallRequests[incomingRequest.Floor][incomingRequest.Button] = true
				newRequests, correctedRequests := hallRequestAssigner(elevID, &globalHallRequests, elevatorMap)

				go distributeRequest(newRequests, fsmChannels.NewRequest)
				go distributeRequest(correctedRequests, fsmChannels.CorrectedRequest)

			}

		case incomingElevator := <-elevSyncChannels.IncomingElevator:
			if incomingElevator != elevatorMap[incomingElevator.ID] {

				updateElevatorMap(elevatorMap, incomingElevator)
				updateGlobalHallRequests(&globalHallRequests, elevatorMap)

				fsmChannels.GlobalHallRequests <- globalHallRequests

			}

		case peersLost := <-elevSyncChannels.PeersLost:
			for _, ID := range peersLost {

				elevator := elevatorMap[ID]
				elevator.Online = false
				elevatorMap[ID] = elevator

				newRequests, correctedRequests := hallRequestAssigner(elevID, &globalHallRequests, elevatorMap)

				elevator.LocalHallRequests = [config.TotalFloors][config.TotalHallButtons]bool{}
				elevatorMap[ID] = elevator

				go distributeRequest(newRequests, fsmChannels.NewRequest)
				go distributeRequest(correctedRequests, fsmChannels.CorrectedRequest)

			}
		case <-elevSyncChannels.BecomeIndependant:
			runIndependantly(elevID, elevatorMap)
		}

	}
}
