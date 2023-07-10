FROM alpine AS worker-runner

RUN apk update && apk upgrade && apk add libc++-dev