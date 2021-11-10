package config

import (

	"github.com/TTK4145-Students-2021/project-group_16/elevio"
)

const (
	TotalFloors      = 4
	TotalButtons     = 3
	TotalHallButtons = 2
)

type ElevBehaviour string

const (
	Idle     ElevBehaviour = "idle"
	Moving   ElevBehaviour = "moving"
	DoorOpen ElevBehaviour = "doorOpen"
	Stop     ElevBehaviour = "stop"
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

type Elevator struct {
	Behaviour         ElevBehaviour         `json:"behaviour"`
	Floor             int                   `json:"floor"`
	Direction         elevio.MotorDirection `json:"direction"`
	CabRequests       [TotalFloors]bool     `json:"cabRequests"`
	LocalHallRequests [TotalFloors][TotalHallButtons]bool
	ID                string
	Obstruction       bool
	Elevstats         Stats
	Online            bool
}

type FsmChannels struct {
	Floor              chan int
	StopBtn            chan bool
	Obstruct           chan bool
	NewRequest         chan elevio.ButtonEvent
	CorrectedRequest   chan elevio.ButtonEvent
	GlobalHallRequests chan [TotalFloors][TotalHallButtons]bool
}


