# Use Go official image
FROM golang:1.24-alpine

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy all source code
COPY . .

# Build the application
RUN go build -o main .

# Expose port (your Gin app runs on 8787)
EXPOSE 8787

# Run the app
CMD ["./main"]
