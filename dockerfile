FROM golang:latest as builder
WORKDIR /app/
COPY . .
RUN go build -o /app/main /app/main.go

FROM alpine:latest
WORKDIR /app/
COPY --from=builder /app/ /app/
ENTRYPOINT ["/app/main"]