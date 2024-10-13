FROM ghcr.io/userver-framework/ubuntu-22.04-userver:latest AS builder
LABEL authors="vikazeni"

COPY ./nb-back /app
WORKDIR /app

RUN adduser user
RUN chown user -R /app

USER user

RUN mkdir build
WORKDIR /app/build
RUN cmake -DCMAKE_BUILD_TYPE=Release ..
RUN cmake --build . --config Release --target nb-back

ENTRYPOINT ["top", "-b"]

FROM ubuntu:22.04
LABEL author="vikazeni"

WORKDIR /app

COPY ./iac/images/nb-back/deps /deps

ARG DEBIAN_FRONTEND=noninteractive

RUN apt update
RUN apt install -y $(cat /deps | tr '\n' ' ')

COPY --from=builder /app/build/nb-back /app/app
COPY --from=builder /app/configs /app/configs

ENTRYPOINT ["./app", "--config", "./configs/static_config.yaml", "--config_vars", "./configs/config_vars.yaml"]