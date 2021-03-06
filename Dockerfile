FROM golang:1.16-alpine AS builder
RUN apk add --no-cache git
WORKDIR /src/pimp
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /bin/pimp

FROM scratch
COPY --from=builder /bin/pimp /bin/pimp
ENTRYPOINT ["/bin/pimp"]
