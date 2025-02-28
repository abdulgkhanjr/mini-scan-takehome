FROM golang:1.20 AS builder

# Build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=1 go build -o processor ./cmd/processor

FROM frolvlad/alpine-glibc
RUN apk update && apk add --no-cache sqlite sqlite-dev
WORKDIR app
COPY --from=builder /src/processor .
EXPOSE 8080
CMD ["/app/processor"]
