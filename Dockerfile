FROM alpine:3.16.2

COPY ./build/linux/jcasc-validator /usr/local/bin
RUN jcasc-validator --version
