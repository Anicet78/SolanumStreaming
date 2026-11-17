FROM oven/bun:slim AS builder

WORKDIR /app

COPY package.json bun.lock* ./

RUN bun install

COPY . .

RUN bun run build

FROM nginx:stable-alpine-slim

WORKDIR /usr/share/nginx

COPY --from=builder /app/dist .

COPY ./nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
