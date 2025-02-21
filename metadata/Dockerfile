FROM golang:1.23-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN go build -o cmd ./metadata/cmd/main.go

#
FROM alpine

RUN mkdir /app

COPY --from=builder /app/cmd /app

CMD [ "/app/cmd" ]