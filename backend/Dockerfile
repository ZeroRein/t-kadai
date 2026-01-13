FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git for fetching dependencies
RUN apk add --no-cache git

COPY . .

# go.sum is not present, so we generate it
ENV GOPROXY=https://proxy.golang.org,direct
RUN go mod tidy
RUN go mod download

RUN go build -o main src/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]
