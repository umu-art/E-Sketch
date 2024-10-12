FROM golang:1.23 as builder
LABEL author="vikazeni"

WORKDIR /app
COPY /nb-proxy .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/app ./src

FROM nginx:stable-alpine3.17
COPY --from=builder /app/build/app /app
ENTRYPOINT ["/app"]