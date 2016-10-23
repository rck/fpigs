FROM alpine:latest
MAINTAINER Roland Kammerer <dev.rck@gmail.com>

# docker run -v $HOME:/fpigs -w /fpigs -it --rm fpigs

ENV FPIGS_VERSION 0.1

RUN apk add --no-cache --virtual .build-deps wget ca-certificates
RUN wget "https://github.com/rck/fpigs/releases/download/v${FPIGS_VERSION}/fpigs-linux-amd64" -O /usr/local/bin/fpigs
RUN chmod +x /usr/local/bin/fpigs
RUN apk del .build-deps

CMD ["fpigs"]
