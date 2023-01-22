FROM golang:1.19.4-alpine3.16 as builder
WORKDIR /usr/src/app
COPY . .
RUN go mod download
RUN go build -o /app ./main.go

FROM alpine:3.16
COPY --from=builder /app .
USER nobody
ENTRYPOINT ["./app"]
