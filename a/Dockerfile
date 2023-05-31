FROM golang:1.19

WORKDIR /service

COPY go.mod go.sum ./
RUN go mod download

COPY a.go ./
COPY ./services/ ./services/
COPY ./handlers/ ./handlers/
COPY ./servicectx/ ./servicectx/

RUN CGO_ENABLED=0 GOOS=linux go build -o /service-a

CMD ["/service-a"]
