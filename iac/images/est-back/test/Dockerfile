FROM ghcr.io/userver-framework/ubuntu-22.04-userver:latest AS builder
LABEL authors="vikazeni"

COPY ./est-back /app
WORKDIR /app

RUN adduser user
RUN chown user -R /app

USER user

RUN mkdir build
WORKDIR /app/build
RUN cmake -DCMAKE_BUILD_TYPE=Release ..
RUN cmake --build . --config Release --target est-back

ENTRYPOINT ["top", "-b"]