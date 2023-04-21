# Use an official Go runtime as a parent image
FROM golang:1.20-alpine AS builder
# Set the working directory to /go/src/app
WORKDIR /usr/src/app
# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify
# Copy the current directory contents into the container at /go/src/app
COPY . ./
# Build the Go program
RUN go build -v -o main ./cmd/main.go

# Use a lightweight Alpine image to run the binary
FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /usr/src/app
COPY --from=builder /usr/src/app/main .
# Expose port 8080 for the application
EXPOSE 8080
# Run the binary program produced by `go build` when the container starts
CMD ["./main"]