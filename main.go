package main

//Main driver for single elevator.
import (
	comm "github.com/TTK4145-Students-2021/project-group_16/communication"
	"github.com/TTK4145-Students-2021/project-group_16/config"
	"github.com/TTK4145-Students-2021/project-group_16/elevio"
	"github.com/TTK4145-Students-2021/project-group_16/fsm"
	"github.com/TTK4145-Students-2021/project-group_16/networkDriver/bcast"
	"github.com/TTK4145-Students-2021/project-group_16/networkDriver/peers"
	oh "github.com/TTK4145-Students-2021/project-group_16/orderHandler"
)

func main() {

	//Port-ID
	elevID := config.GetPortID()

	//Ports
	//peersPort := 15010 //To update active elevators
	bcastPort := 15020
	orderPort := 15030 //To assign a new order to an elevator

	//Channels
	fsmChannels := config.FsmChannels{
		Floor:    make(chan int),
		StopBtn:  make(chan bool),
		Obstruct: make(chan bool),
		NewOrder: make(chan elevio.ButtonEvent),
	}
	//Communication channel
	commChannels := config.CommChannels{
		TransmitEnable:   make(chan bool),
		PeerUpdate:       make(chan peers.PeerUpdate),
		UpdateElevator:   make(chan config.Elevator),
		OutgoingElevator: make(chan config.Elevator),
		IncomingElevator: make(chan config.Elevator),
		OutgoingOrder:    make(chan elevio.ButtonEvent),
		IncomingOrder:    make(chan elevio.ButtonEvent),
	}
	//OrderHandlerChannel

	ohChannels := config.OrderHandlerChannels{

		IncomingOrder:    make(chan elevio.ButtonEvent),
		IncomingElevator: make(chan config.Elevator),
	}

	//Lokal channel for sensoring new orders.
	var (
		newOrder = make(chan elevio.ButtonEvent)
	)

	//Threads

	//Every new order are broadcasted and given to cost function.
	go elevio.PollButtons(newOrder)
	go elevio.PollFloorSensor(fsmChannels.Floor)
	go elevio.PollObstructionSwitch(fsmChannels.Obstruct)
	go elevio.PollStopButton(fsmChannels.StopBtn)

	go fsm.SingleElevator(elevID, fsmChannels, commChannels)

	go comm.Communication(elevID, newOrder, commChannels, ohChannels, fsmChannels)

	go oh.OrderHandler(elevID, commChannels, ohChannels, fsmChannels)

	//go peers.Transmitter(peersPort, elevID, commChannels.TransmitEnable)
	//go peers.Receiver(peersPort, commChannels.PeerUpdate)

	go bcast.Transmitter(orderPort, commChannels.OutgoingOrder)
	go bcast.Receiver(orderPort, commChannels.IncomingOrder)

	go bcast.Transmitter(bcastPort, commChannels.OutgoingElevator)
	go bcast.Receiver(bcastPort, commChannels.IncomingElevator)

	select {}
}
