cd C:\Users\super\Documents\ITU\Semestre\3rd\Distributed\Assignment5
start cmd /k "cd Server & go run server.go 0"
start cmd /k "cd Server & go run server.go 1"
start cmd /k "cd Server & go run server.go 2"
start cmd /k "cd Server & go run server.go 3"

start cmd /k "cd Client & go run client.go 0 3"
start cmd /k "cd Client & go run client.go 1 3"
start cmd /k "cd Client & go run client.go 2 3"