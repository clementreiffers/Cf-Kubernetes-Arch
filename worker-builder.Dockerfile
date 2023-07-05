FROM ubuntu AS worker-builder

RUN apt-get update && apt-get install -y clang libc++-dev nodejs npm