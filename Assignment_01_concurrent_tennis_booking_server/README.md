### Compilation flags:
```bash
go build client/tennis_client.go
go build server/tennis_server.go
```

The binaries attached (in dir `build`) have been compiled statically.

### Usage:
```bash
./tennis_client infile.csv host:port
./tennis_server port > outfile.csv
```

### Assumptions Made:
- There are 3 queues for each game type: Single, Double, Both/Any
- Preference is given to the first player that arrives,  
specifically the least timestamp of the players in front of all 3 queues is given priority in game type.  
Eg: (Queue{Timestamp}) `Single {5}, Double {10}, Both {4b}` will result in a Single game as `4b` prefers Single.
- If the preferred game type is not feasible at the moment, then the other game type is checked.
- Requests may come in any manner, `ArrivalTime` mentioned in the request may be lower than the one already in queue,  
but the request will be added to the back of the queue.
- The Scheduling Output is done at the server side, nothing is returned to the client.
- All scheuling is controlled by a single mutex, so it is coarse-grained.