# Use the official Golang image as the base image
FROM golang AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application source code into the container
COPY . .

# Build the Go application
# RUN go build -o main ./user-service/main.go
RUN go build -o main main.go

# Set execute permissions for the wait-for-port.sh script
# RUN chmod +x wait-for-port.sh

# Second stage: Create a smaller image
FROM alpine

# Set the working directory inside the container
WORKDIR /app

# Copy only the binary and necessary files from the builder stage
# COPY --from=builder /app/main .
# COPY wait-for-port.sh  ./

# Expose the port the application runs on
EXPOSE 8080

# Command to run the application
CMD ["/app"]






# FROM golang

# RUN mkdir /app

# ADD . /app

# WORKDIR /app

# RUN go build -o main main.go

# EXPOSE 8080
# CMD ["/app/main"]
