FROM node:alpine AS worker-builder

RUN apk update && apk upgrade && apk add clang libc++-dev

RUN yarn global add workerd