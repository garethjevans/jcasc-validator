FROM alpine:3.18.2

COPY ./build/linux/jcasc-validator /usr/local/bin
RUN jcasc-validator --version
