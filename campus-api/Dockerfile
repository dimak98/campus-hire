FROM golang:1.22 AS builder

WORKDIR /srv

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o campus-api ./cmd/api/main.go

FROM alpine:3.18.3

WORKDIR /srv/

COPY --from=builder /srv/campus-api .

ENTRYPOINT ["./campus-api"]
