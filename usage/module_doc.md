
** Overview elevio module **

A server acts like the elevator. By starting the simulator we can give call and commands by pressing buttons given in simulator.con file. 

When we start the simulator we can express ports, commands etc in the con file. 

- fsm.go 

- elevator_io.go 
Includes poll, get, set functions to sensor the elevator. Has an init function to connect an elevator to network. Basically the hardware module from last year. 

- requests.go 

- timer.go 
Timer module, start and stop timer. using "time" 