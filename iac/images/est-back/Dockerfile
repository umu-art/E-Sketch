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
RUN git clone https://github.com/open-telemetry/opentelemetry-cpp.git

WORKDIR /opentelemetry-cpp
RUN git submodule update --init --recursive
RUN mkdir build

WORKDIR /opentelemetry-cpp/build
RUN cmake -DBUILD_TESTING=OFF -DBUILD_SHARED_LIBS=ON -DWITH_OTLP_HTTP=ON -DOPENTELEMETRY_INSTALL=ON ..
RUN make && make install

# Copy source code
COPY ./est-back /app/est-back
COPY ./api/build/est-back-cpp /app/api/build/est-back-cpp

# Build source code
WORKDIR /app/est-back/build
RUN cmake ..
RUN make

# TODO: remake imgae

ENTRYPOINT ["./est-back"]