FROM golang:latest as builder
ENV GOPROXY="https://goproxy.cn"
ENV CGO_ENABLED=0
WORKDIR /BBingyan/
COPY go.mod .
COPY go.sum .
RUN go mod tidy && go mod verify
COPY . .
RUN go build -o /build/app ./cmd/main.go

FROM alpine:latest
WORKDIR /usr/bin/BBingyan
COPY --from=builder /build/app ./app
COPY --from=builder /BBingyan/Config ./Config
