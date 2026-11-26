FROM oven/bun:slim AS builder

WORKDIR /app

COPY package.json bun.lock* ./

RUN bun install

COPY . .

ARG VITE_API_URL
ARG VITE_AUTH_PATH
ARG VITE_MOVIES_PATH
ARG VITE_STREAM_PATH

ENV VITE_API_URL=$VITE_API_URL
ENV VITE_AUTH_PATH=$VITE_AUTH_PATH
ENV VITE_MOVIES_PATH=$VITE_MOVIES_PATH
ENV VITE_STREAM_PATH=$VITE_STREAM_PATH

RUN bun run build

FROM nginx:stable-alpine-slim

WORKDIR /usr/share/nginx

COPY --from=builder /app/dist .

COPY ./nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
