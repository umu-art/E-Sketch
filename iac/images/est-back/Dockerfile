FROM ghcr.io/userver-framework/ubuntu-22.04-userver:latest AS builder
LABEL authors="vikazeni"

COPY ./est-back /app/est-back
COPY ./api/build/est-back-cpp /app/api/build/est-back-cpp
WORKDIR /app/est-back

RUN adduser user
RUN chown user -R /app

USER user

RUN mkdir build
WORKDIR /app/est-back/build
RUN cmake -DCMAKE_BUILD_TYPE=Release ..
RUN cmake --build . --config Release --target est-back

FROM ubuntu:22.04
LABEL author="vikazeni"

WORKDIR /app

COPY ./iac/images/est-back/deps /deps

ARG DEBIAN_FRONTEND=noninteractive

RUN apt update
RUN apt install -y $(cat /deps | tr '\n' ' ')

COPY --from=builder /app/est-back/build/est-back /app/app
COPY --from=builder /app/est-back/configs /app/configs

ENTRYPOINT ["./app", "--config", "./configs/static_config.yaml", "--config_vars", "./configs/config_vars.yaml"]