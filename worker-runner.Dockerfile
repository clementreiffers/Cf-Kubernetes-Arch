FROM debian:bullseye AS worker-runner

RUN apt-get update && apt-get install -y libc++-dev