FROM golang:alpine3.15 AS build

WORKDIR /app

COPY app/go.mod ./
COPY app/go.sum ./

RUN go mod download

COPY app/ ./

RUN go build -o /webapp

EXPOSE 8080

CMD ["/webapp"]