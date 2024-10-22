FROM drogonframework/drogon AS builder
LABEL authors="vikazeni"

RUN apt-get update
RUN apt-get install -y postgresql-all

COPY ./est-back /app/est-back
COPY ./api/build/est-back-cpp /app/api/build/est-back-cpp

WORKDIR /app/est-back/build
RUN cmake ..
RUN make

ENTRYPOINT ["top", "-b"]