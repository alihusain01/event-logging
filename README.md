# MP0 - Event Logging

## Names:
Ali Husain (alijh2), Satej Sukthankar (satejrs2)

## Cluster Numbers:
4201 - 4210

## Git URL:
https://gitlab.engr.illinois.edu/cs4251732105/mp0-event-logging/-/tree/main

## Git Revision Number:
1e42f4a3f90ffb3167930cf3e5da24539241d579

# Getting started

## Starting server
1. ``` cd server ``` to enter the server
2. ``` go run server.go {server port number} ``` to start the centralized logger

## Starting node
1. ``` cd node ``` to enter the node folder
2. ``` make build ``` to compile the node.go file
3. ``` python3 -u generator.py {frequency} | ./node {node_name} {server IP} {server port number} ``` to start the node events

## Conducting Analytics
1. Upon terminating the server, a ```logFile.txt``` folder will be created

