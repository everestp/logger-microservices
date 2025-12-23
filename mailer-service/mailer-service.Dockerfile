# Stage 1: Build the binary
FROM golang:1.24-alpine3.21 AS builder

RUN mkdir /app
COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 go build -o mailerApp ./cmd/api
RUN chmod +x /app/mailerApp

# Build a tiny docker image
FROM alpine:latest

RUN mkdir /app

# Copy the binary
COPY --from=builder /app/mailerApp /app/mailerApp

# --- ADD THIS LINE ---
# This copies your templates folder into the /app directory in the container
COPY ./templates /app/templates 

WORKDIR /app

EXPOSE 8083

CMD [ "/app/mailerApp" ]