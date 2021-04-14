
### Documentation for driver module go

Driveren snakker med server. 

* Global variables
- bool _initialized 
- int _numFloors 
- sync.Mutex _mtx
- _conn net.Conn
These global variables are needed for using network module. 

* Enumerations :
- MotorDirection int (1,-1,0)
- ButtonType int (0,1,2)

* Structs: 
ButtonEvent {
    Floor int 
    Button Buttontype
} 
// Denne inneholder orderen representert med et heltall, og button type


* Init 
Makes a dial connection with string "localhost:port", and int number of floors. initialises sync object. Initialize true. 

** Get/Set functions ** 
These functions are bounded to the server

Get functions uses mutexes to secure only one thread is accessing the shared variable. 

* getButton (button ButtonType, floor int) bool 
Denne skriver til serveren et array med 6 bytes, button og floor. Også returnerer den true om den button er trykka på. 

* getfloor()
Writes to server the floor-message. Reads feedback. Returns current floor if the elevator is on a floor, otherwise -1. 

* getStop()
Writes stop message to server. Reads feedback. Returns bool whether we have stopped or not. 

* getObstruction()
Writes obstruction message to server. Reads feedback. Returns bool whether we have obstruction or not. 

* SetMotorDirection(dir MotorDirection)
Writes to server, set motor actuation dir. 

* SetButtonLamp(button ButtonType, floor int, value bool)
Writes to server, set button lamp

* SetFloorIndicator(floor int)
Writes to server, set floor indicator

* SetDoorOpenLamp(value bool)
Writes to server, set lamp

* SetStopLamp(value bool)
Writes to server, set stop lamp 

** Poll functions **
We use poll functions for an event. The for-select-channel statement wil handle every event. 

* PollButtons (chan<- ButtonEvent)

Makes an order matrix where orders not being on same floor is pushed into receiver channel. 

* PollObstructuonSwitch(chan<- bool)



* Message types: 
- []byte{7,0,0,0} get current floor
- []byte{8,0,0,0} get stop
- []byte{9,0,0,0} get obstruction
