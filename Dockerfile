FROM alpine:3.18.0

COPY ./build/linux/jcasc-validator /usr/local/bin
RUN jcasc-validator --version
