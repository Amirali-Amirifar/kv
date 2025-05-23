# kv
Kv is a distributed key-value database, with high availability and eventual consistency 

## Using the client with a single node setup
```sh 
go run ./cmd/kvNode 
go run ./cmd/kvClient http://localhost:8081
KV Database Client
Successfully connected to host http://localhost:8081

Type HELP for available commands or QUIT to exit

kv> set "user:1" "John Doe" 
OK
kv> get "user:1"
"John Doe"
kv> 



```