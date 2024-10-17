FROM nginx:stable-alpine3.17
LABEL author="vikazeni"

COPY /nb-front/build/ /usr/share/nginx/html
COPY /iac/images/nb-front/nginx.conf /etc/nginx/conf.d/default.conf