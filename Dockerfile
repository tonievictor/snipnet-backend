FROM golang:1.22.5-bullseye AS builder
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o nest

FROM alpine:3.20.2
RUN apk --no-cache add ca-certificates libc6-compat
WORKDIR /app 
COPY --from=builder /usr/src/app/nest .
CMD ["/app/nest"]
