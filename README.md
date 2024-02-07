# MP0 - Event Logging

## Names:
Ali Husain (alijh2), Satej Sukthankar (satejrs2)

# Getting started

## Starting server
1. ``` cd server ``` to enter the server
2. ``` go run server.go {server port number} ``` to start the centralized logger

## Starting node
1. ``` cd node ``` to enter the node folder
2. ``` make build ``` to compile the node.go file
3. ``` python3 -u generator.py {frequency} | ./node {node_name} {server IP} {server port number} ``` to start the node events

## Conducting Analytics
Ensure you have ```matplotlib```, ```numpy```, and ```pandas``` installed before running analytics.

1. Upon terminating the server, a logFile.txt file will be created in the server folder.
2. Enter the analytics folder and run the code blocks within each notebook for their corresponding analytic.
3. A graph will be generated within the analytics folder for the corresponding notebook.


