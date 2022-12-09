FROM golang:1.19.3-alpine as build

WORKDIR /go/build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o buybot ./app

FROM alpine:3.17 as release
WORKDIR /app

COPY --from=build /go/build/ ./

CMD ["/app/buybot"]