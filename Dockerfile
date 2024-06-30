FROM golang:1.22.4-bullseye

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

WORKDIR /usr/src/app/cmd

RUN go build -v -o /usr/local/bin/nest ./...

CMD ["nest"]
