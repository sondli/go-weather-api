FROM golang:1.21 as builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . ./
RUN go build -o server

FROM debian:latest as production
COPY --from=builder /app/server /app/server
CMD ["/app/server"]
