# Causal Broadcast Simulator

This is a program that simulates delivery of messages between n processes, using vector clocks. This is done to guarantee
causal delivery of messages (i.e. no causal anomalies).  

The program introduces random artificial delays in broadcasting the messages to other processes, to simulate delays in the network. The causal delivery of messages is then demonstrated, by comparing the vector clocks of the incoming message and that of the current process. Messages that do not violate causal delivery are immediately delivered, and messages from the "causal future" are buffered to be delivered later.  

When a message is delivered, the process checks the buffer if any buffered message(s) are ready for delivery, based on the updated vector clock value.

## Assumptions

- Message sends are the only events (message receives/internal events are not considered, but can be easily added)
- The underlying TCP sockets that are used for inter-process communication delivers the messages in order.

## How to Run

- Clone the repo  
```
git clone https://github.com/krithikvaidya/distsys-implementations.git
```

- Change directory  
```
cd distsys-implementations/causal_broadcast
```

- Decide the number of processes involved (n). Based on this, open (n - 1) more terminal windows and run each process (default value of n is 3):  
```
go run . -n <value chose for n>
```

- Follow the on-screen instructions

## TODO:

- Improve buffer checking efficiency

## Reference

[UCSC CSE138 Lectures](https://www.youtube.com/watch?v=G0wpsacaYpE&list=PLNPUF5QyWU8O0Wd8QDh9KaM1ggsxspJ31)