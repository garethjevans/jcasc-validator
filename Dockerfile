FROM alpine:3.17.1

COPY ./build/linux/jcasc-validator /usr/local/bin
RUN jcasc-validator --version
