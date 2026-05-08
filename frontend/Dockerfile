FROM node:slim

WORKDIR /app

COPY ./nginx.conf /etc/nginx/conf.d/default.conf:ro

COPY . .

RUN npm run build

FROM nginx:stable-alpine-slim

WORKDIR /usr/share/nginx

COPY --from=builder /app/build .

EXPOSE 8080

ENTRYPOINT [ "nginx -g 'daemon off;'" ]
