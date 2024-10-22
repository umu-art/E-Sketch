FROM ubuntu:24.10 AS builder
LABEL authors="vikazeni"

RUN apt-get update
RUN apt-get install -y postgresql-all git gcc g++ cmake libjsoncpp-dev uuid-dev zlib1g-dev

WORKDIR /
RUN git clone https://github.com/drogonframework/drogon

WORKDIR /drogon
RUN git submodule update --init
RUN mkdir build

WORKDIR /drogon/build
RUN cmake ..
RUN make && make install

COPY ./est-back /app/est-back
COPY ./api/build/est-back-cpp /app/api/build/est-back-cpp

WORKDIR /app/est-back/build
RUN cmake ..
RUN make

FROM ubuntu:24.10

RUN apt-get update
RUN apt-get install -y libpq-dev libjsoncpp-dev uuid-dev

COPY --from=builder /app/est-back/build/est-back /app/app
COPY --from=builder /app/est-back/config.json /config.json

WORKDIR /app

ENTRYPOINT ["./app"]