FROM golang:1.17.0-alpine as builder

# Install git.
RUN apk update && apk add --no-cache git

# Set the current working directory inside the container 
WORKDIR /app

# Copy go mod and sum files 
COPY go.mod go.sum ./

# Download all dependencies.
RUN go mod download 

# Copy the source from the current directory to the working Directory inside the container 
COPY . .

# Build the Go app
RUN go build  -o main .

# Expose port 8080 
EXPOSE 8080

#Command to run the executable
CMD ["./main"]
