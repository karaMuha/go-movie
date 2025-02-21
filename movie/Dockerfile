FROM golang:1.23-alpine as builder

RUN mkdir /app
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN go build -o cmd ./movie/cmd/main.go

#
FROM alpine

RUN mkdir /app

COPY --from=builder /app/cmd /app

CMD [ "/app/cmd" ]