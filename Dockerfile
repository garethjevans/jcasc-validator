FROM alpine:3.17.0

COPY ./build/linux/jcasc-validator /usr/local/bin
RUN jcasc-validator --version
