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

Subscriber to all messages:
```
curl -N -X GET -H 'Accept: text/event-stream' http://localhost:4242/rpc/Chat/SubscribeMessages
```

## Go client
Send message:
```
go run ./client-go Hello there!
```

Subscriber to all messages:
```
go run ./client-go
```
