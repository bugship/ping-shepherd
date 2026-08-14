FROM golang:1.25-alpine AS build
WORKDIR /src
RUN apk add --no-cache ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/ping-shepherd ./cmd/shepherd

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=build /out/ping-shepherd /usr/local/bin/ping-shepherd
USER nobody
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/ping-shepherd"]
