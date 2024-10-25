FROM ubuntu:24.10 AS builder
LABEL authors="vikazeni"

RUN apt-get update

# Install Drogon
RUN apt-get install -y libpq-dev git gcc g++ cmake libjsoncpp-dev uuid-dev zlib1g-dev

WORKDIR /
RUN git clone https://github.com/drogonframework/drogon

WORKDIR /drogon
RUN git submodule update --init
RUN mkdir build

WORKDIR /drogon/build
RUN cmake ..
RUN make && make install

# Install OpenTelemetry
RUN apt-get install -y protobuf-compiler curl libcurl4-openssl-dev

WORKDIR /
RUN git clone https://github.com/abseil/abseil-cpp.git

WORKDIR /abseil-cpp/build
RUN cmake ..
RUN make && make install

WORKDIR /
RUN git clone https://github.com/open-telemetry/opentelemetry-cpp.git

WORKDIR /opentelemetry-cpp/build
RUN cmake -DBUILD_TESTING=OFF -DWITH_OTLP_HTTP=ON -DWITH_ABSEIL=ON -DOPENTELEMETRY_INSTALL=ON ..
RUN make && make install
