FROM alpine:latest

COPY ./etc/timezone /etc/timezone

COPY ./etc/localtime /etc/localtime

COPY ./main /bin/kk-trip

RUN chmod +x /bin/kk-trip

COPY ./config /config

COPY ./app.ini /app.ini

ENV KK_ENV_CONFIG /config/env.ini

VOLUME /config

CMD kk-trip $KK_ENV_CONFIG

