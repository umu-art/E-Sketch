FROM ghcr.io/userver-framework/ubuntu-22.04-userver:latest AS builder
LABEL authors="vikazeni"

COPY ./est-back /app/est-back
COPY ./api/build/est-back-cpp /app/api/build/est-back-cpp
WORKDIR /app/est-back

RUN adduser user
RUN chown user -R /app

USER user

RUN make test-debug
ENTRYPOINT ["echo", "completed"]