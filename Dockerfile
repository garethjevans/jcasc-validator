FROM alpine:3.14.3

COPY ./build/linux/jcasc-validator /usr/local/bin
RUN jcasc-validator --version
