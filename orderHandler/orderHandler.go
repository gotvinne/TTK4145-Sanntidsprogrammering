package orderHandler

import (
	"fmt"

	"github.com/TTK4145-Students-2021/project-group_16/config"
	"github.com/TTK4145-Students-2021/project-group_16/elevio"
)

//Make distribute function, update

func OrderHandler(elevID string, commCh config.CommChannels, ohChannels config.OrderHandlerChannels, fsmChannels config.FsmChannels) {

	//Initializing elevator table
	elevatorTable := make(map[string]config.Elevator)
	//Initializing hall requests, default value
	hallRequests := [config.TotalFloors][config.TotalHallButtons]bool{}

	for {
		select {

		case incomingOrder := <-ohChannels.IncomingOrder:

			if incomingOrder.Button == elevio.BTCab { //Only if local
				fsmChannels.NewOrder <- incomingOrder
			} else {
				onlineElevators := GetOnlineElevators(elevatorTable)

				//Clear hallRequests due to logic in synchronizeHallRequest. One elevator never clears out Hall Order.
				hallRequests = [config.TotalFloors][config.TotalHallButtons]bool{}
				UpdateHallRequests(&hallRequests, onlineElevators)

				//Add Hallorder in hallRequests:
				hallRequests[incomingOrder.Floor][incomingOrder.Button] = true
				fmt.Println(hallRequests)
				costOrder := CostOrder{HallRequests: hallRequests, OnlineElevators: onlineElevators}

				assignedOrders := CostFunction(costOrder)
				fmt.Println(assignedOrders)

				localAssignedOrders := assignedOrders[elevID]
				localHallRequests := elevatorTable[elevID].HallRequests

				//Check if elevator locally has a new order
				newOrder := AssignedNewOrder(localAssignedOrders, localHallRequests)
				if newOrder.Floor != -1 {
					fsmChannels.NewOrder <- newOrder
				}
			}

		case incomingElevator := <-ohChannels.IncomingElevator:
			UpdateElevatorTable(elevatorTable, incomingElevator)
			//Denne kjører 4 ganger
			UpdateHallRequests(&hallRequests, elevatorTable)

			//Update lights locally:
			go func() {

			}()

		}

	}
}
