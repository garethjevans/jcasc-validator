FROM alpine:3.15.4

COPY ./build/linux/jcasc-validator /usr/local/bin
RUN jcasc-validator --version
