FROM alpine:3.13.5

COPY ./build/linux/jcasc-validator /usr/local/bin
RUN jcasc-validator --version
