FROM golang:1.19.1 AS build

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY ./ ./
RUN go build -o weather-bot-discord ./

FROM alpine:latest as release
COPY --from=build /usr/src/app/weather-bot-discord ./
COPY zip-codes.csv ./
RUN apk add --no-cache libc6-compat # fix for go executables
RUN apk add --no-cache tzdata
CMD ["./weather-bot-discord"]