package communication

import (
	"fmt"
	"time"

	"github.com/TTK4145-Students-2021/project-group_16/config"
	elevSync "github.com/TTK4145-Students-2021/project-group_16/elevatorSynchronizer"
	"github.com/TTK4145-Students-2021/project-group_16/elevio"
	"github.com/TTK4145-Students-2021/project-group_16/networkDriver/peers"
)

const (
	broadCastInteval = 15 * time.Millisecond
)

type GlobalRequest struct {
	Request   elevio.ButtonEvent
	TimeStamp time.Time
}

type CommChannels struct {
	PeersTransmitEnable chan bool
	PeerUpdate          chan peers.PeerUpdate

	UpdateElevator       chan config.Elevator
	UpdateOnlineElevator chan string

	OutgoingElevator chan config.Elevator
	IncomingElevator chan config.Elevator
	OutgoingRequest  chan GlobalRequest
	IncomingRequest  chan GlobalRequest
}

func Communication(elevID string, newLocalRequest chan elevio.ButtonEvent, commChannels CommChannels, elevSyncChannels elevSync.ElevSyncChannels) {

	currentElevator := config.Elevator{}
	currentRequest := GlobalRequest{}

	broadCastTicker := time.NewTicker(broadCastInteval)

	for {
		select {
		case newLocalRequest := <-newLocalRequest:

			elevSyncChannels.IncomingRequest <- newLocalRequest

			if newLocalRequest.Button != elevio.BTCab {
				newGlobalRequest := GlobalRequest{Request: newLocalRequest, TimeStamp: time.Now()}
				commChannels.OutgoingRequest <- newGlobalRequest
				currentRequest = newGlobalRequest
			}

		case incommingRequest := <-commChannels.IncomingRequest:
			if incommingRequest.TimeStamp.After(currentRequest.TimeStamp) {
				elevSyncChannels.IncomingRequest <- incommingRequest.Request
				currentRequest = incommingRequest
			}

		case updateElevator := <-commChannels.UpdateElevator:

			elevSyncChannels.IncomingElevator <- updateElevator

			if updateElevator.Online {
				commChannels.OutgoingElevator <- updateElevator
			}

			currentElevator = updateElevator

		case incomingElevator := <-commChannels.IncomingElevator:
			if incomingElevator.Online {
				elevSyncChannels.IncomingElevator <- incomingElevator
			}
		case peerUpdate := <-commChannels.PeerUpdate:

			if peerUpdate.New != "" {
				commChannels.OutgoingElevator <- currentElevator
			}
			if len(peerUpdate.Lost) == 1 {
				elevSyncChannels.PeersLost <- peerUpdate.Lost
			} else if len(peerUpdate.Lost) > 1 {

				elevSyncChannels.BecomeIndependant <- true
			}
			fmt.Println("\n\nPeers: ", peerUpdate.Peers)
			fmt.Println("New: ", peerUpdate.New)
			fmt.Println("Lost: ", peerUpdate.Lost)
			fmt.Println("\n\n")

		case <-broadCastTicker.C:
			commChannels.OutgoingRequest <- currentRequest
			commChannels.OutgoingElevator <- currentElevator

		}
	}
}
