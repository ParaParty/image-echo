FROM docker.io/golang:alpine3.16 AS build

WORKDIR /work
COPY . .
RUN export GO111MODULE=on && export GOPROXY=https://goproxy.cn && go build -o main

FROM docker.io/alpine:3.16.3

WORKDIR /work
COPY --from=build /work/main .
ENTRYPOINT ["./main"]