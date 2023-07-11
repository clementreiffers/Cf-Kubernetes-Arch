FROM node:bullseye AS worker-builder

RUN apt-get update
RUN apt-get install -y libc++-dev
RUN apt-get install -y clang
RUN yarn global add workerd

