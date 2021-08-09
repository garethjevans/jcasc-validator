FROM alpine:3.14.1

COPY ./build/linux/jcasc-validator /usr/local/bin
RUN jcasc-validator --version
