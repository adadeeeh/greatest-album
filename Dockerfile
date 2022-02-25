FROM golang:alpine3.15 AS build

WORKDIR /app

COPY app/go.mod ./
COPY app/go.sum ./

RUN go mod download

COPY app/ ./

RUN go build -o webapp

# EXPOSE 8080

# CMD ["/app/webapp"]


FROM bash:4.0.44

WORKDIR /app

COPY --from=build /app/.env ./
COPY --from=build /app/webapp ./

EXPOSE 8080

ENTRYPOINT [ "/app/webapp" ]