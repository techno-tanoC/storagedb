FROM golang:1.16 AS build
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY main.go storage.go ./
RUN CGO_ENABLED=1 go build -trimpath -ldflags "-s -w -buildid="


FROM debian:10-slim
WORKDIR /app

RUN apt update && apt install -y ca-certificates

COPY --from=build /build/storagedb /app/storagedb

CMD [ "/app/storagedb" ]
