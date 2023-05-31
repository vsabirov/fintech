FROM golang:1.19

WORKDIR /service

COPY go.mod go.sum ./
RUN go mod download

COPY b.go ./
COPY ./consumers/ ./consumers/
COPY ./migrations/ ./migrations/
COPY ./servicectx/ ./servicectx/
COPY ./services/ ./services/

RUN CGO_ENABLED=0 GOOS=linux go build -o /service-b

CMD ["/service-b"]
