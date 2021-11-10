package main

import (
	comm "github.com/TTK4145-Students-2021/project-group_16/communication"
	"github.com/TTK4145-Students-2021/project-group_16/config"
	elevSync "github.com/TTK4145-Students-2021/project-group_16/elevatorSynchronizer"
	"github.com/TTK4145-Students-2021/project-group_16/elevio"
	"github.com/TTK4145-Students-2021/project-group_16/fsm"
	"github.com/TTK4145-Students-2021/project-group_16/networkDriver/bcast"
	"github.com/TTK4145-Students-2021/project-group_16/networkDriver/peers"
)

const (
	peersPort   = 15010
	bcastPort   = 15020
	requestPort = 15030
)

func main() {

	elevID := fsm.GetPortID()

	fsmChannels := config.FsmChannels{
		Floor:              make(chan int),
		Obstruct:           make(chan bool),
		NewRequest:         make(chan elevio.ButtonEvent, 8),
		CorrectedRequest:   make(chan elevio.ButtonEvent, 8),
		GlobalHallRequests: make(chan [config.TotalFloors][config.TotalHallButtons]bool, 100),
	}
	commChannels := comm.CommChannels{
		PeersTransmitEnable:  make(chan bool),
		PeerUpdate:           make(chan peers.PeerUpdate),
		UpdateElevator:       make(chan config.Elevator),
		UpdateOnlineElevator: make(chan string),
		OutgoingElevator:     make(chan config.Elevator),
		IncomingElevator:     make(chan config.Elevator),
		OutgoingRequest:      make(chan comm.GlobalRequest, 8),
		IncomingRequest:      make(chan comm.GlobalRequest, 8),
	}

	elevSyncChannels := elevSync.ElevSyncChannels{

		IncomingRequest:   make(chan elevio.ButtonEvent, 8),
		IncomingElevator:  make(chan config.Elevator),
		PeersLost:         make(chan []string),
		BecomeIndependant: make(chan bool),
	}

	var (
		buttonEvent = make(chan elevio.ButtonEvent)
	)

	go elevio.PollButtons(buttonEvent)
	go elevio.PollFloorSensor(fsmChannels.Floor)
	go elevio.PollObstructionSwitch(fsmChannels.Obstruct)

	go fsm.RunElevator(elevID, fsmChannels, commChannels)

	go comm.Communication(elevID, buttonEvent, commChannels, elevSyncChannels)

	go elevSync.ElevatorSynchronizer(elevID, elevSyncChannels, fsmChannels)

	go peers.Transmitter(peersPort, elevID, commChannels.PeersTransmitEnable)
	go peers.Receiver(peersPort, commChannels.PeerUpdate)

	go bcast.Transmitter(requestPort, commChannels.OutgoingRequest)
	go bcast.Receiver(requestPort, commChannels.IncomingRequest)

	go bcast.Transmitter(bcastPort, commChannels.OutgoingElevator)
	go bcast.Receiver(bcastPort, commChannels.IncomingElevator)

	select {}
}
