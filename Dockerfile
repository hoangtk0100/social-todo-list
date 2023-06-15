# Build stage
FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o todo-app .

# Run stage
FROM alpine
WORKDIR /app
COPY --from=builder /app/todo-app .
COPY app.env .
COPY db.env .
COPY start.sh .
COPY wait-for.sh .
RUN chmod +x start.sh wait-for.sh

EXPOSE 9090
ENTRYPOINT ["/app/start.sh"]
CMD ["/app/todo-app"]
