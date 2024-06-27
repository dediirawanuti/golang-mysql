# Start with the official Golang image
FROM golang:1.22

# Set the current working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code from the scripts directory
COPY scripts/ /app/scripts/

# Set the current working directory inside the container
WORKDIR /app/scripts

# Create the images directory
RUN mkdir -p /app/images

# Build the Go app
RUN go build -o main .

# Expose the port on which the app will run
EXPOSE 8910

# Command to run the executable
CMD ["./main"]
