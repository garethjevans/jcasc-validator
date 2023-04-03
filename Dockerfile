FROM alpine:3.17.3

COPY ./build/linux/jcasc-validator /usr/local/bin
RUN jcasc-validator --version
