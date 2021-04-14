package config

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/TTK4145-Students-2021/project-group_16/elevio"
	"github.com/TTK4145-Students-2021/project-group_16/networkDriver/peers"
)

const (
	TotalFloors      = 4
	TotalElev        = 3
	TotalButtons     = 3
	TotalHallButtons = 2
)

//fsm
type ElevBehaviour string

const (
	Idle     ElevBehaviour = "idle"
	Moving   ElevBehaviour = "moving"
	DoorOpen ElevBehaviour = "doorOpen"
	Stop     ElevBehaviour = "stop"
)

const (
	DoorOpenTime = 2 * time.Second
)

type ClearRequestVariant int

const (
	// Assume everyone waiting for the elevator gets on the elevator, even if
	// they will be traveling in the "wrong" direction for a while
	CVAll ClearRequestVariant = iota

	// Assume that only those that want to travel in the current direction
	// enter the elevator, and keep waiting outside otherwise
	CVInDirn
)

type Stats struct {
	CRV              ClearRequestVariant
	DoorOpenDuration float64
}

type FsmChannels struct {
	Floor    chan int
	StopBtn  chan bool
	Obstruct chan bool
	NewOrder chan elevio.ButtonEvent
}

type Elevator struct {
	Behaviour    ElevBehaviour                       `json:"behaviour"`
	Floor        int                                 `json:"floor"`
	Direction    elevio.MotorDirection               `json:"direction"`
	CabRequests  [TotalFloors]bool                   `json:"cabRequests"`
	HallRequests [TotalFloors][TotalHallButtons]bool `json:"hallRequests"`

	ID          string
	Obstruction bool
	Stop        bool
	Elevstats   Stats
	Online      bool
}

//Order handler

type OrderHandlerChannels struct {
	IncomingOrder    chan elevio.ButtonEvent
	IncomingElevator chan Elevator
}

//Connection

type CommChannels struct {
	TransmitEnable chan bool
	PeerUpdate     chan peers.PeerUpdate

	UpdateElevator   chan Elevator
	OutgoingElevator chan Elevator
	IncomingElevator chan Elevator
	OutgoingOrder    chan elevio.ButtonEvent
	IncomingOrder    chan elevio.ButtonEvent
}

func GetPortID() string {

	elevIDPtr := flag.Int("id", 0, "HeisID")
	flag.Parse()

	if *elevIDPtr < 15000 {
		fmt.Println("Missing or wrongly typed portID")
		os.Exit(1)
	}

	elevID := strconv.Itoa(*elevIDPtr)
	return elevID
}

func StateToString(State ElevBehaviour) string {
	switch State {
	case Idle:
		return "Idle"
	case DoorOpen:
		return "DoorOpen"
	case Moving:
		return "Moving"
	case Stop:
		return "Stop"
	}
	return "Invalid State"
}

func DirToString(Direction elevio.MotorDirection) string {
	switch Direction {
	case elevio.MDUp:
		return "MDUp"
	case elevio.MDDown:
		return "MDDown"
	case elevio.MDStop:
		return "MDStop"
	}
	return "Invalid Direction"
}

func PrintElevator(e Elevator) {
	fmt.Printf("\n")
	fmt.Printf("\n")

	fmt.Printf("ID: ")
	fmt.Printf("%v", e.ID)
	fmt.Printf("\n")

	fmt.Printf("State: ")
	fmt.Printf("%v", StateToString(e.Behaviour))
	fmt.Printf("\n")

	fmt.Printf("Floor: ")
	fmt.Printf("%v", e.Floor)
	fmt.Printf("\n")

	fmt.Printf("Direction: ")
	fmt.Printf("%v", DirToString(e.Direction))
	fmt.Printf("\n")

	fmt.Printf("Table: ")
	fmt.Printf("\n")

	for btn := 0; btn < 2; btn++ {
		for floor := 0; floor < TotalFloors; floor++ {
			fmt.Printf("%v", e.HallRequests[floor][btn])
			fmt.Printf("     ")
		}
		fmt.Printf("\n")
	}
	for floor := 0; floor < TotalFloors; floor++ {
		fmt.Printf("%v", e.CabRequests[floor])
		fmt.Printf("     ")
	}
	fmt.Printf("\n")
}

func PrintButtonEvent(order elevio.ButtonEvent) {
	fmt.Println("\n")
	fmt.Println("Floor:")
	fmt.Println(order.Floor)
	fmt.Println("Buttontype")
	fmt.Println(reflect.ValueOf(order.Button).Interface().(elevio.ButtonType))

}
