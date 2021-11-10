**FSM** 
--------------------------------------------------

The FSM module represents a finite state machine implementation of a single elevator. This module recieves hall- and cabrequests from elevatorSynchronizer and signals from the elevio module. Moreover, the FSM module is responsible for:

- Initializing the elevator
- Executing hall- and cabrequests
- Storing requests in backup files
- Keeping door open when the obstruction signal is high
- Setting the door timer and handle its timeout
- Setting the motor error timer and handle its timeout
- Setting lights based on global hallrequests
- Sending its local state to the communication module
