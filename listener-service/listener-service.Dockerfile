# Stage 1: Build the binary
# Changed to alpine3.21 which supports newer Go versions
FROM golang:1.24-alpine3.21 AS builder

RUN mkdir /app
COPY . /app

WORKDIR /app

# This will now succeed because the Go version matches or exceeds your go.mod
RUN CGO_ENABLED=0 go build -o listenerApp .
RUN chmod +x /app/listenerApp

# Build a tiny docker image
FROM alpine:latest

# Fixed: Added the space between mkdir and /app
RUN mkdir /app

COPY --from=builder /app/listenerApp /app/listenerApp

# Best practice: Inform Docker about the port
EXPOSE 8084

CMD [ "/app/listenerApp" ]