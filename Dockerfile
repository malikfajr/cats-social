FROM golang:1.22.1-alpine AS builder
WORKDIR app
COPY . .
RUN go get
RUN go build -o main main.go

FROM alpine
WORKDIR app

COPY --from=builder /app/main .

EXPOSE 8080

ENV DB_NAME=cats-social
ENV DB_PORT=5432
ENV DB_HOST=localhost
ENV DB_USERNAME=postgres
ENV DB_PASSWORD=secret
ENV DB_PARAMS=sslmode=disable
ENV JWT_SECRET=VERY_s3cr3t
ENV BCRYPT_SALT=8

CMD ["/app/main"]