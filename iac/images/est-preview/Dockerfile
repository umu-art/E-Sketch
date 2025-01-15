FROM node:20
LABEL author="vikazeni"

WORKDIR /usr/src/app
COPY /est-preview/dist ./dist

ENTRYPOINT ["node", "dist/main"]