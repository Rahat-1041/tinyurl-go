FROM golang:1.21-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o tinyurl main.go

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/tinyurl .
EXPOSE 8080
CMD ["./tinyurl"]
