#!/bin/sh
dmd -w -g src/sim_server.d src/timer_event.d -ofSimElevatorServer
osascript -e 'tell application "System Events" to tell process "Terminal" to keystroke "t" using command down'
osascript -e 'tell app "Terminal" to do script "cd Simulator-v2-1.5; ./SimElevatorServer --port 15000" in window 1'
osascript -e 'tell application "System Events" to tell process "Terminal" to keystroke "t" using command down'
osascript -e 'tell app "Terminal" to do script "go build; go run main.go --id 15000" in window 1'

osascript -e 'tell application "System Events" to tell process "Terminal" to keystroke "t" using command down'
osascript -e 'tell app "Terminal" to do script "cd Simulator-v2-1.5; ./SimElevatorServer --port 15001" in window 1'
osascript -e 'tell application "System Events" to tell process "Terminal" to keystroke "t" using command down'
osascript -e 'tell app "Terminal" to do script "go build; go run main.go --id 15001" in window 1'

osascript -e 'tell application "System Events" to tell process "Terminal" to keystroke "t" using command down'
osascript -e 'tell app "Terminal" to do script "cd Simulator-v2-1.5; ./SimElevatorServer --port 15002" in window 1'
go build
go run main.go --id 15002