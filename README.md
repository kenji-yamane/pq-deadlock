# P-out-of-Q model deadlock

Implementation of the Kshemkalyani-Singhal Algorithm for P-out-of-Q Model on a simulation of a distributed system with ipc

## Compilation instructions

The only needed dependency in this project
is go version 1.17.

To build all targets:
```bash
make build-all
```

## Execution instructions

To execute a simulation of three different processes, first instantiate
the processes simulators:
```bash
make p1-3
make p2-3
make p3-3
```
On three different terminals.

Then, use the comands to send REQUEST and REPLY messages on each terminal to design a WFG  
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

## Fixtures

We also have an alternative way to build WFG, using
a csv file. It is read by the `puppeteer` target, which
then applies the requested commands to each
process. The following is an example of a valid csv file
```csv
1,ask 1 2,3
2,ask 1 3,3
3,ask 1 1,3
2,liberate 1,3
```
The first column is the id of the process to which the command
is directed, the second is the argument in itself, and the third
is the number of milliseconds you want the program to wait after
executing the command.

