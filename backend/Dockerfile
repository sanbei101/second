FROM golang:1.26-alpine as builder
WORKDIR /app
COPY go.mod go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOEXPERIMENT=jsonv2 go build -o server .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .

ENTRYPOINT ["/app/server"]

EXPOSE 8848