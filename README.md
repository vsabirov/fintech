# Services
Service A - RESTful HTTP Server, listens for the following endpoints:
```
POST { id: <string>, amount: <number>, receiver: <string> } -> /accounts/<account-id>/transfer -> 201 or 404
```

Sends all requests to service B through kafka.
Go + Echo + Kafka-Go


Service B - Request operator, receives requests from kafka and sends responses to kafka.
Go + Kafka-Go + PGX

# Service dependencies
Service A depends on Kafka

Service B depends on Kafka and Postgres.

Service A and Service B don't depend on eachother. If Service A stops, the REST API will stop working and new requests won't be received. If Service B stops, handling of Kafka messages will be suspended until Service B gets back online.