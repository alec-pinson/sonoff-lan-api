FROM golang:1.20-alpine AS build_base
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN cd cmd/sonoff-lan-api && \
    go build -o /build/out/my-app .

# Start fresh from a smaller image
FROM alpine:3.17.2
RUN apk add ca-certificates openssl
COPY --from=build_base /build/out/my-app /app/sonoff-lan-api
WORKDIR /app
CMD ["/app/sonoff-lan-api"]
