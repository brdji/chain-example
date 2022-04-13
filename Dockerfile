#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
RUN apk add build-base
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
# build will fail if output dir is not specified beforehand
RUN mkdir ./build
RUN go build -o ./build -v ./...

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN mkdir /app
COPY --from=builder /go/src/app/build /app
WORKDIR /app
RUN chmod +x /app/chain-listener
ENTRYPOINT /app/chain-listener
LABEL Name=chainlistener Version=0.0.1
EXPOSE 8081
