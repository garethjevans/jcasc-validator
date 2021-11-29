FROM alpine:3.15.0

COPY ./build/linux/jcasc-validator /usr/local/bin
RUN jcasc-validator --version
