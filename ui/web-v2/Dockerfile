FROM nginx:1.13.11-alpine

RUN apk update \
  && apk add ca-certificates \
  && update-ca-certificates \
  && ln -sf /var/nginx/nginx.conf /etc/nginx/nginx.conf

WORKDIR /etc/nginx

COPY dist/apps/admin/ /var/www/

CMD ["nginx", "-g", "daemon off;"]

EXPOSE 80 443
