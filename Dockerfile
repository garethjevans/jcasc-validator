FROM alpine:3.20.1

COPY ./build/linux/jcasc-validator /usr/local/bin
RUN jcasc-validator --version
