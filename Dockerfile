# Use an official Go runtime as a parent image
FROM golang:alpine

# Set the working directory to /go/src/app
WORKDIR /go/src/app

# Copy the current directory contents into the container at /go/src/app
COPY . ./

# Build the Go program
RUN go build -o ./app ./cmd/main.go

# Expose port 8080 for the application
EXPOSE 8080

# Run the Go program when the container starts
CMD ["./app"]
