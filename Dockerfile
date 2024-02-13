FROM golang:1.22 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ cmd/
COPY pkg/ pkg/

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:3.19.1 as release

WORKDIR /

COPY --from=build /app/main .

EXPOSE 8080

ENTRYPOINT ["/main"]
