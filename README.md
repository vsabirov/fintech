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

# Usage
Run
```
docker compose up -d
```
in the root directory of the project. Wait until all containers are built and started.


You can interact with Kafka and Postgres using kowl and pgweb (available at ports 8080 and 8081 respectively), and send HTTP requests to Service A using any HTTP client (Insomnia, for example). Service A will be available at port 2500.

Create 2 users with any ids and balance total through pgweb. 
![create-accounts](/docs/images/create-accounts.png)
![list-accounts](/docs/images/list-accounts.png)

Send a HTTP request to Service A.
![create-transfer](/docs/images/create-transfer.png)

Service A logs:
```
time="2023-05-31T18:38:58+03:00" level=info msg=request URI=/accounts/a/transfer status=201
```

Service B logs:
```
time="2023-05-31T18:38:58+03:00" level=info msg="Processing new message." message="{\"id\":\"1\",\"amount\":500.5,\"receiver\":\"b\",\"sender\":\"a\"}" topic=transfer

time="2023-05-31T18:38:58+03:00" level=info msg="Transfer request processed successfully." message="{\"id\":\"1\",\"amount\":500.5,\"receiver\":\"b\",\"sender\":\"a\"}" topic=transfer
```

View new account totals.
![new-accounts-total](/docs/images/new-accounts-total.png)

View new registered transfer.
![transfer](/docs/images/transfer.png)

Try to create a transfer with the same id.

Service B logs:
```
time="2023-05-31T18:43:28+03:00" level=info msg="Processing new message." message="{\"id\":\"1\",\"amount\":500.5,\"receiver\":\"b\",\"sender\":\"a\"}" topic=transfer

time="2023-05-31T18:43:28+03:00" level=error msg="Failed to process transfer request." error="Transfer with this ID already exists." message="{\"id\":\"1\",\"amount\":500.5,\"receiver\":\"b\",\"sender\":\"a\"}" topic=transfe
```

No new transfer was created.

Transfers can be randomly marked as invalid after 30 seconds from processing. In that case, the transfer will be removed from database and account funds will return back to their owner.

```
2023-05-31 20:12:34 time="2023-05-31T17:12:34Z" level=warning msg="Transfer was marked invalid, trying to restore account funds." request="{2 500.5 b a}"

2023-05-31 20:12:34 time="2023-05-31T17:12:34Z" level=info msg="Funds restored successfully." request="{2 500.5 b a}"
```