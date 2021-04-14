package communication

import (
	"github.com/TTK4145-Students-2021/project-group_16/config"
	"github.com/TTK4145-Students-2021/project-group_16/elevio"
)

func Communication(elevID string, NewOrder chan elevio.ButtonEvent, commChannels config.CommChannels, ohChannels config.OrderHandlerChannels, fsmChannels config.FsmChannels) {
	//TransmitEnable tell peers if online

	//commChannel.TransmitEnable <- true
	for {
		select {

		// //case peerUpdate := <-commChannel.PeerUpdate:
		// 	fmt.Printf("%v", "\n")
		// 	fmt.Printf("%v", "\n")
		// 	fmt.Printf("%v", "Peers: ")
		// 	fmt.Printf("%v", peerUpdate.Peers)
		// 	fmt.Printf("%v", "\n")
		// 	fmt.Printf("%v", "New elevator: "+peerUpdate.New+"\n")
		// 	fmt.Printf("%v", "Peers lost: ")
		// 	fmt.Printf("%v", peerUpdate.Lost)
		// 	fmt.Printf("%v", "\n")

		//Order distribution
		case newLocalBtnEvent := <-NewOrder:
			//Distribute locally
			ohChannels.IncomingOrder <- newLocalBtnEvent

			//Distribute Globally, to thread in main
			if newLocalBtnEvent.Button != elevio.BTCab {
				commChannels.OutgoingOrder <- newLocalBtnEvent
			}
		case incommingOrder := <-commChannels.IncomingOrder:
			ohChannels.IncomingOrder <- incommingOrder

		case updateElevator := <-commChannels.UpdateElevator:
			//Distrubute locally:
			ohChannels.IncomingElevator <- updateElevator

			//Distribute Globally:
			commChannels.OutgoingElevator <- updateElevator

		case incomingElevator := <-commChannels.IncomingElevator:
			ohChannels.IncomingElevator <- incomingElevator
			//Broadcast order

			//case updatePeerOrders := <-commChannel.UpdatePeerOrderTable:
			//fmt.Println(updatePeerOrders)
			//case e := <-commChannel.Update:
			/*
				case peerUpdate := <-commChannels.PeerUpdate:
				 	fmt.Printf("%v", "\n")
				 	fmt.Printf("%v", "\n")
				 	fmt.Printf("%v", "Peers: ")
				 	fmt.Printf("%v", peerUpdate.Peers)
				 	fmt.Printf("%v", "\n")
				 	fmt.Printf("%v", "New: "+peerUpdate.New+"\n")
				 	fmt.Printf("%v", "Lost: ")
				 	fmt.Printf("%v", peerUpdate.Lost)
				 	fmt.Printf("%v", "\n")
			*/
		}
	}
}
