package elevatorSynchronizer

import (
	"encoding/json"
	"os/exec"

	"github.com/TTK4145-Students-2021/project-group_16/config"
	"github.com/TTK4145-Students-2021/project-group_16/elevio"
)

func updateGlobalHallRequests(hallRequests *[config.TotalFloors][config.TotalHallButtons]bool, elevMap map[string]config.Elevator) {
	for _, elevator := range elevMap {
		synchronizeGlobalHallRequests(hallRequests, elevator)
	}
}

func synchronizeGlobalHallRequests(hallRequests *[config.TotalFloors][config.TotalHallButtons]bool, elev config.Elevator) {
	for floor := 0; floor < config.TotalFloors; floor++ {
		for btn := 0; btn < config.TotalHallButtons; btn++ {
			if elev.LocalHallRequests[floor][btn] {
				hallRequests[floor][btn] = true
			}
		}
	}
}

func updateElevatorMap(elevMap map[string]config.Elevator, elev config.Elevator) {
	if elev.Online {
		elevMap[elev.ID] = elev
	}
}

func getOnlineElevators(elevMap map[string]config.Elevator) map[string]config.Elevator {

	onlineElevators := make(map[string]config.Elevator)

	for _, elevator := range elevMap {
		if elevator.Online {
			onlineElevators[elevator.ID] = elevator
		}
	}
	return onlineElevators
}

func costFunction(costRequest CostRequest) map[string][config.TotalFloors][config.TotalHallButtons]bool {
	formatted, err := json.Marshal(costRequest)
	if err != nil {
		panic(err)
	}
	output, err := exec.Command("./hall_request_assigner", "--input", string(formatted)).Output()
	assignedRequests := make(map[string][config.TotalFloors][config.TotalHallButtons]bool)
	_ = json.Unmarshal(output, &assignedRequests)

	return assignedRequests
}

func assignedNewRequests(assignedHallRequests [config.TotalFloors][config.TotalHallButtons]bool, localHallRequests [config.TotalFloors][config.TotalHallButtons]bool) ([]elevio.ButtonEvent, []elevio.ButtonEvent) {
	var newRequestTable []elevio.ButtonEvent
	var correctedRequestTable []elevio.ButtonEvent
	for floor := 0; floor < config.TotalFloors; floor++ {
		for btn := 0; btn < config.TotalHallButtons; btn++ {
			if assignedHallRequests[floor][btn] == true && localHallRequests[floor][btn] == false {
				//If VHS = true and RHS = false: Create newRequest
				newRequest := elevio.ButtonEvent{Floor: floor, Button: elevio.ButtonType(btn)}
				newRequestTable = append(newRequestTable, newRequest)
			} else if assignedHallRequests[floor][btn] == false && localHallRequests[floor][btn] == true {
				//If VHS = false and RHS = true: Create correctedRequest
				correctedRequest := elevio.ButtonEvent{Floor: floor, Button: elevio.ButtonType(btn)}
				correctedRequestTable = append(correctedRequestTable, correctedRequest)
			}
		}
	}
	if len(newRequestTable) == 0 {
		newRequestTable = append(newRequestTable, elevio.ButtonEvent{Floor: -1, Button: elevio.ButtonType(0)})
	} else if len(correctedRequestTable) == 0 {
		correctedRequestTable = append(correctedRequestTable, elevio.ButtonEvent{Floor: -1, Button: elevio.ButtonType(0)})
	}
	return newRequestTable, correctedRequestTable
}

func hallRequestAssigner(elevID string, hallRequests *[config.TotalFloors][config.TotalHallButtons]bool, elevMap map[string]config.Elevator) ([]elevio.ButtonEvent, []elevio.ButtonEvent) {
	onlineElevators := getOnlineElevators(elevMap)
	updateGlobalHallRequests(hallRequests, elevMap)

	costRequest := CostRequest{HallRequests: *hallRequests, OnlineElevators: onlineElevators}

	assignedHallRequestMap := costFunction(costRequest)

	localAssignedHallRequests := assignedHallRequestMap[elevID]
	localHallRequests := elevMap[elevID].LocalHallRequests
	newAssignedRequestTable, correctedRequestTable := assignedNewRequests(localAssignedHallRequests, localHallRequests)

	return newAssignedRequestTable, correctedRequestTable
}

func runIndependantly(elevID string, elevMap map[string]config.Elevator) {

	for ID, elevator := range elevMap {
		if ID != elevID {
			elevator.Online = false
			elevator.LocalHallRequests = [config.TotalFloors][config.TotalHallButtons]bool{}
			elevMap[ID] = elevator
		}
	}
}

func distributeRequest(requests []elevio.ButtonEvent, channel chan<- elevio.ButtonEvent) {

	for _, request := range requests {
		if request.Floor != -1 {
			channel <- request
		}
	}
}
