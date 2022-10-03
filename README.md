# WebRPC SSE discovery playground

## Run Go server at http://localhost:4242 with:
```bash
go run ./server-go
```

## CURL client
Send message:
```
curl -X POST -H 'Content-Type: application/json' --data '{"author": "Vojtech", "msg": "Hello there!"}' http://localhost:4242/rpc/Chat/SendMessage
```

Subscribe to all messages:
```
curl -N -X GET -H 'Accept: text/event-stream' http://localhost:4242/rpc/Chat/SubscribeMessages
```

## Go client
```
go run ./client-go
```

Joins chat.
Send a message by typing and hitting enter.