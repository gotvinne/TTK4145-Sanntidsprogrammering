package orderHandler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/TTK4145-Students-2021/project-group_16/config"
	"github.com/TTK4145-Students-2021/project-group_16/elevio"
)

type CostOrder struct {
	HallRequests    [config.TotalFloors][config.TotalHallButtons]bool `json:"hallRequests"`
	OnlineElevators map[string]config.Elevator                        `json:"states"`
}

//For test
type id struct {
	Behaviour   string  `json:"behaviour"`
	Floor       int     `json:"floor"`
	Direction   string  `json:"direction"`
	CabRequests [4]bool `json:"cabRequests"`
}

type test struct {
	HallRequests [4][2]bool  `json:"hallRequests"`
	States       interface{} `json:"states"`
}

func UpdateHallRequests(hallRequests *[config.TotalFloors][config.TotalHallButtons]bool, et map[string]config.Elevator) {
	for _, elevator := range et {
		SynchronizeHallRequests(hallRequests, elevator)
	}
}
func SynchronizeHallRequests(hallRequests *[config.TotalFloors][config.TotalHallButtons]bool, e config.Elevator) {

	for floor := 0; floor < config.TotalFloors; floor++ {
		for btn := 0; btn < config.TotalHallButtons; btn++ {
			if e.HallRequests[floor][btn] {
				hallRequests[floor][btn] = true
			}
		}
	}
}

//Slices/arrays are passed by reference in golang.
func UpdateElevatorTable(et map[string]config.Elevator, e config.Elevator) {
	if e.Online {
		et[e.ID] = e
	}

	//PrintElevatorTable(et)
}

//Funka
func GetOnlineElevators(et map[string]config.Elevator) map[string]config.Elevator {

	onlineElevators := make(map[string]config.Elevator)

	for _, elevator := range et {
		if elevator.Online {
			onlineElevators[elevator.ID] = elevator
		}
	}
	return onlineElevators
}

//Update with costorder
//chmod +x hall_request_assigner
func CostFunction(costOrder CostOrder) map[string][config.TotalFloors][config.TotalHallButtons]bool {
	formatted, err := json.Marshal(costOrder)
	if err != nil {
		panic(err)
	}
	output, err := exec.Command("./hall_request_assigner", "--input", string(formatted)).Output()
	assignedOrders := make(map[string][config.TotalFloors][config.TotalHallButtons]bool)
	_ = json.Unmarshal(output, &assignedOrders)

	return assignedOrders
}

func AssignedNewOrder(assignedHallRequests [config.TotalFloors][config.TotalHallButtons]bool, localHallRequests [config.TotalFloors][config.TotalHallButtons]bool) elevio.ButtonEvent {

	for floor := 0; floor < config.TotalFloors; floor++ {
		for btn := 0; btn < config.TotalHallButtons; btn++ {
			if assignedHallRequests[floor][btn] != localHallRequests[floor][btn] {
				newOrder := elevio.ButtonEvent{Floor: floor, Button: elevio.ButtonType(btn)}

				return newOrder
			}
		}
	}

	return elevio.ButtonEvent{Floor: -1, Button: elevio.ButtonType(0)}

}

func TestCostFunction() [4][2]bool {

	one := id{
		Behaviour:   "moving",
		Floor:       2,
		Direction:   "up",
		CabRequests: [4]bool{false, false, true, true},
	}

	two := id{
		Behaviour:   "idle",
		Floor:       1,
		Direction:   "stop",
		CabRequests: [4]bool{false, false, false, false},
	}
	type rar struct {
		Id_1 id `json:"id_1"`
		Id_2 id `json:"id_2"`
	}
	rarere := rar{
		Id_1: one,
		Id_2: two,
	}

	test := test{
		HallRequests: [4][2]bool{{false, false}, {true, false}, {false, false}, {false, true}},
		States:       rarere,
	}

	bytes, err := json.Marshal(test)
	if err != nil {
		panic(err)
	}

	var data map[string]interface{}
	output, err := exec.Command("./hall_request_assigner", "--input", string(bytes)).Output()
	fmt.Println(string(output))
	_ = json.Unmarshal(output, &data)
	data_1 := data["id_1"].([]interface{})

	//Save in file
	file, _ := json.MarshalIndent(data_1, "", " ")
	_ = ioutil.WriteFile("test69.json", file, 0644)

	var correctElev [4][2]bool

	_ = json.Unmarshal([]byte(file), &correctElev)
	return correctElev
}

func PrintAllElevatorTable(et map[string]config.Elevator) {

	if et == nil {
		fmt.Println("Map is empty")
	}
	for _, elevator := range et {
		config.PrintElevator(elevator)
	}
}

func PrintLocalElevatorTable(ID string, et map[string]config.Elevator) {

	if et == nil {
		fmt.Println("Map is empty")
	}
	for _, elevator := range et {
		if elevator.ID == ID {
			config.PrintElevator(elevator)
		}
	}
}

func SetAllHallLights(et map[string]config.Elevator) {

	for _, elevator := range et {

	}

}
