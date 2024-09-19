# syntax=docker/dockerfile:1
FROM busybox:latest
WORKDIR /app

ARG goarch=amd64

ADD https://github.com/chuccp/smtp2http/releases/latest/download/smtp2http-linux-${goarch}.tar.gz /app/
RUN tar -xzf smtp2http-linux-${goarch}.tar.gz && rm -rf  *.tar.gz && chmod a+x /app/smtp2http
EXPOSE 12566 12567
CMD [ "/app/smtp2http","-web_port","12566","-api_port","12577"  ]