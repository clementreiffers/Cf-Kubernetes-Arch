FROM ubuntu AS worker

RUN apt-get update && apt-get install -y libc++-dev