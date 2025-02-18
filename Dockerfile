FROM golang:1.23

LABEL maintainer="Ariel Llanos"

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

ENV PORT 8080

RUN make build

CMD ["./minesweeper-api"]
