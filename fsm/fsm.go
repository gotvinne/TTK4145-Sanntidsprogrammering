package fsm

import (
	"time"

	comm "github.com/TTK4145-Students-2021/project-group_16/communication"
	"github.com/TTK4145-Students-2021/project-group_16/config"
)

const (
	doorOpenDuration    = 3 * time.Second
	timeUntilMotorError = 10 * time.Second
)

type BackupRequests struct {
	BackupCabRequests       [config.TotalFloors]bool                          `json:"backUpCabRequests"`
	BackupLocalHallRequests [config.TotalFloors][config.TotalHallButtons]bool `json:"backupLocalHallRequests"`
}


func RunElevator(elevID string, fsmChannels config.FsmChannels, commChannels comm.CommChannels) {

	doorTimer := time.NewTimer(doorOpenDuration)
	doorTimer.Stop()

	motorErrorTimer := time.NewTimer(timeUntilMotorError)
	motorErrorTimer.Stop()

	elev := initializeElevator(elevID, fsmChannels.Floor, doorTimer, motorErrorTimer)

	for {
		select {
		case currentFloor := <-fsmChannels.Floor:
			arrivedAtFloor(&elev, currentFloor, doorTimer, motorErrorTimer)
			commChannels.PeersTransmitEnable <- true
			elev.Online = true
			saveRequestsToBackup(elev)
			commChannels.UpdateElevator <- elev

		case newRequest := <-fsmChannels.NewRequest:
			newRequestHandler(&elev, newRequest, doorTimer, motorErrorTimer)
			saveRequestsToBackup(elev)
			commChannels.UpdateElevator <- elev

		case correctedRequest := <-fsmChannels.CorrectedRequest:
			localHallRequests := elev.LocalHallRequests
			localHallRequests[correctedRequest.Floor][int(correctedRequest.Button)] = false
			elev.LocalHallRequests = localHallRequests
			saveRequestsToBackup(elev)
			commChannels.UpdateElevator <- elev

		case obstructionHigh := <-fsmChannels.Obstruct:
			obstructionHandler(&elev, obstructionHigh, doorTimer)
			commChannels.UpdateElevator <- elev

		case <-doorTimer.C:
			onDoorTimeout(&elev, doorTimer, motorErrorTimer, commChannels.PeersTransmitEnable)
			commChannels.UpdateElevator <- elev

		case <-motorErrorTimer.C:
			motorErrorHandler(&elev, commChannels.PeersTransmitEnable, motorErrorTimer)
			commChannels.UpdateElevator <- elev

		case globalRequests := <-fsmChannels.GlobalHallRequests:
			setAllHallLights(globalRequests)
		}
	}
}
