FROM nginx:stable-alpine3.17
LABEL author="vikazeni"

COPY /est-front/build/ /usr/share/nginx/html
COPY /iac/images/est-front/nginx.conf /etc/nginx/conf.d/default.conf