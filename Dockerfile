# Build stage
FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o todoApp .

# Run stage
FROM alpine
WORKDIR /app
COPY --from=builder /app/todoApp .
EXPOSE 9090
ENTRYPOINT ["./todoApp"]
