FROM golang:1.23-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o payflow ./cmd/server
FROM alpine:3.20
WORKDIR /app
COPY --from=build /app/payflow .
EXPOSE 8080
CMD ["./payflow"]
