FROM golang:1.22.7-alpine as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /apod ./cmd/main.go

FROM alpine:3.20.3

ENV SHUTDOWN_SERVER_TIMEOUT_SEC=10

WORKDIR /app 
COPY --from=builder /apod /app/
# swagger documentation
COPY docs /app/docs/

ENTRYPOINT ["/app/apod" ]
