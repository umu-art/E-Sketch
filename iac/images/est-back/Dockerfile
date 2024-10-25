FROM vikazeni/est-back-build-container:latest AS builder
LABEL authors="vikazeni"

# Copy source code
COPY ./est-back /app/est-back
COPY ./api/build/est-back-cpp /app/api/build/est-back-cpp

# Build source code
WORKDIR /app/est-back/build
RUN cmake ..
RUN make

# TODO: remake imgae

ENTRYPOINT ["./est-back"]