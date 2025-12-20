# Stage 1: Build the binary
FROM golang:1.24-alpine3.21 AS builder

RUN mkdir /app
# Copy the entire service source code
COPY . /app

WORKDIR /app


# Compile the Go binary
RUN CGO_ENABLED=0 go build -o authApp ./cmd/api

# Ensure binary is executable
RUN chmod +x /app/authApp

# Stage 2: Build tiny runtime image
FROM alpine:latest

# Fixed: Added the space between mkdir and /app
RUN mkdir /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/authApp /app/authApp

# Expose the port your service listens on
EXPOSE 8081

# Run the binary
CMD ["/app/authApp"]
