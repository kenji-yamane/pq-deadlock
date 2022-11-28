# P-out-of-Q model deadlock

Implementation of the Kshemkalyani-Singhal Algorithm for P-out-of-Q Model on a simulation of a distributed system with ipc

## Instructions

To build all targets:
```bash
make build-all
```

Then, to execute a simulation of three different processes, first instantiate
the processes simulators:
```bash
make p1-3
make p2-3
make p3-3
```
On three different terminals.

Then, use the comands for send REQUEST and REPLY messages on each terminal for design a WFG  
REQUEST:  
```bash
ask p_number list_of_nodes
```
REPLY:
```bash
liberate parent_number
```

To detect a deadlock, choose a terminal, then:
```bash
detect
```

In case you want to build your own simulation with a different
number of processes, you can do so executing the binary directly.
For instance, to simulate four different processes:
```bash
./bin/process 1 10005 10004 10003 10002
./bin/process 2 10005 10004 10003 10002
./bin/process 3 10005 10004 10003 10002
./bin/process 4 10005 10004 10003 10002
```
The first argument indicates to the simulator which is his own port,
from the arguments that follows.
