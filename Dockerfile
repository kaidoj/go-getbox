FROM golang:alpine AS build-env
ENV GO111MODULE=on
WORKDIR /go/src/getbox

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o main

# final stage
FROM alpine
WORKDIR /app
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=build-env /go/src/getbox/config.yml .
COPY --from=build-env /go/src/getbox/main . 
ENTRYPOINT ./main