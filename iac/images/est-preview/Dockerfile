FROM node:20-alpine
LABEL author="vikazeni"

WORKDIR /usr/src/app
COPY /est-preview/ ./est-preview
COPY /lib/ ./lib

WORKDIR /usr/src/app/est-preview

ENTRYPOINT ["node", "dist/main"]