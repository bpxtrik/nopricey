FROM golang:1.27 AS build

WORKDIR /src

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/nopricey ./cmd

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /usr/local/go/lib/time/zoneinfo.zip /zoneinfo.zip
COPY --from=build /out/nopricey /app/nopricey

ENV ZONEINFO=/zoneinfo.zip

WORKDIR /data

EXPOSE 10000

ENTRYPOINT ["/app/nopricey"]
