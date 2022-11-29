# How to run program
1. Open an arbitrary number of terminals
2. Set the path of some of the terminals to the project and enter the server folder
3. Type ```'go run server.go x'``` in as many terminals as you want servers
4. Enter the client folder instead of the server folder
5. Type ```'go run client.go x y'``` in as many terminals as you want clients,
    ```x``` = The id of the server/client
    ```y``` = Total amount of servers created
    
    NOTE: The x-values have to start from 0 and be consecutive. Here is an example:
```go
    go run server.go 0
    go run server.go 1

    go run client.go 0 2
    go run client.go 1 2
    go run client.go 2 2
```
When the servers and clients are created, you can use the client terminals to either write a number that is going to be your bid, or you can write 'result' to see the state of the auction.

All operations can be seen in the corresponding log files for each server/client in the client/server folder

#### NOTE: Optionally if you are using Windows, you may also be able to use one of the script files to start the program for you, with a default of 4 servers and 3 clients.
- WINDOWS: change the path of the startup.bat file and double click the file