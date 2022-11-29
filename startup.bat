cd *PATH GOES HERE*
start cmd /k "cd Server & go run server.go 0"
start cmd /k "cd Server & go run server.go 1"
start cmd /k "cd Server & go run server.go 2"
start cmd /k "cd Server & go run server.go 3"

start cmd /k "cd Client & go run client.go 0 4"
start cmd /k "cd Client & go run client.go 1 4"
start cmd /k "cd Client & go run client.go 2 4"