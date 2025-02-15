FROM golang:1.23-alpine AS builder

WORKDIR /build



ADD go.mod .
ADD go.sum .
RUN go mod download
RUN pwd


COPY cmd cmd
COPY internal internal

RUN GOOS=linux go build -o app ./cmd

FROM golang:1.23-alpine
WORKDIR /root/

ENV POSTGRES_USER=${POSTGRES_USER}
ENV POSTGRES_PASSWORD=${POSTGRES_PASSWORD}

COPY --from=builder /build/app .

EXPOSE 8080

CMD ["./app"]