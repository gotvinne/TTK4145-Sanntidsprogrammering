**Communication module** 
--------------------------------------------------
The communication module has following responsibilities: 
- Distribution of a occured hall requests 
- Distribution of an update in a local elevator 
- Managing packet loss 
- Handling online/offline states on network

Due to the functionality of the given network driver, our system utilises a peer-to-peer network between elevators. In addition messages are broadcasted by json formatting using the go library "json". This implies that every message broadcasted on network is received by all peers. The communication module receives event from newBtnEvent channel and elevator information on updateElevator. After broadcasting the communication module receives information in the network on the incomming channels, thus the information is sent to the elevatorSynchronizer.