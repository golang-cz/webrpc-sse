# WebRPC SSE discovery playground

## Run Go server at http://localhost:4242 with:
```bash
go run ./server-go
```

## CURL client
```
curl -X POST -H 'Content-Type: application/json' --data '{"author": "Vojtech", "msg": "Hello there!"}' http://localhost:4242/rpc/Chatbot/SendMessage
```

```
curl -N -X GET -H 'Content-Type: text/event-stream' http://localhost:4242/rpc/Chatbot/SubscribeMessages
```
