# Services
Service A - RESTful HTTP Server, listens for the following endpoints:
```
GET /accounts/<account-id>/balance -> { total: <number> } or 404
POST { id: <string>, amount: <number>, receiver: <string> } -> /accounts/<account-id>/transfer -> 201 or 404
```

Sends all requests to service B through kafka.
Go + Echo + Kafka-Go


Service B - Request operator, receives requests from kafka and sends responses to kafka.
Go + Kafka-Go