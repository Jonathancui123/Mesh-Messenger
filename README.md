# go-mesh-network
![Video demo of go-mesh-network!](https://github.com/Jonathancui123/go-mesh-network/blob/master/demo/go-mesh-network%20demo.gif)

## About
A simple mesh network client that sends messages over TCP connections. It records the addresses of all peers it has encountered, and also broadcasts incoming messages to existing peers. New clients can join the network by dialing into any current member of the network. 

- Built with goRoutines that wait for input, listen for incoming messages, and send messages to peers

- Duplicated messages due to broadcasting are ignored by tracking unique ID's for each message

- To populate the list of peers, messages contain the client's current address


## Setup
Firewall must be opened up to incoming and outbound connections.

```
go install go-mesh-network
go-mesh-network [ -nickname <your-nickname> -dial <peer's host:port to dial>] 
```
