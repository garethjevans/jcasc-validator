FROM alpine:3.19.0

COPY ./build/linux/jcasc-validator /usr/local/bin
RUN jcasc-validator --version
