# WebRPC SSE discovery playground

## Run Go server at http://localhost:4242 with:
```bash
go run ./server-go
```

## CURL client
```
curl -X POST -H 'Origin: http://localhost:4242' -H 'Content-Type: application/json' --data '{"author": "Vojtech", "msg": "Hello there!"}' http://localhost:4242/rpc/Chat/SendMessage
```

```
curl -N -X GET -H 'Origin: http://localhost:4242' -H 'Accept: text/event-stream' http://localhost:4242/rpc/Chat/SubscribeMessages
```
