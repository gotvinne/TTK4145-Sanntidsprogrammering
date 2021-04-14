
Ressurser: https://github.com/TTK4145/Project-resources

Go-Driver: https://github.com/TTK4145-Students-2021/driver-go 

Network-Driver: https://github.com/TTK4145-Students-2021/Network-go/blob/master/main.go

Message struct:

Meldinger som skal sendes:


Arrived at floor

Check-ups: 
- Do we need Keypress (send over ID with newOrder)


Known missbehaviour: 
- First time when initializing a new elevator, it takes the next order anyway due to missing update from elevator: Solution: Make peers module such that if a new elevator is recognized send elevator table. 

- If three elevators are on three distinct floors, and a hall order is given, everyone opens a door even tho floor is wrong. 
