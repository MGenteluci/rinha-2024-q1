FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

ADD cmd ./cmd
ADD pkg ./pkg

RUN CGO_ENABLED=0 GOOS=linux go build cmd/main.go

EXPOSE 8080

CMD ["./main"]
