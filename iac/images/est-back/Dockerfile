FROM drogonframework/drogon AS builder
LABEL authors="vikazeni"

COPY ./est-back /app/est-back
COPY ./api/build/est-back-cpp /app/api/build/est-back-cpp
WORKDIR /app/est-back



ENTRYPOINT ["top", "-b"]