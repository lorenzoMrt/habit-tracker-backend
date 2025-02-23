# Use the official Golang image as the base image
FROM golang:1.23.4-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-s -w" -o main ./cmd/api


# Start a new stage from scratch
FROM scratch
COPY --from=builder /app/main /main
EXPOSE 8080
CMD ["/main"]
