FROM golang:1.25.8 AS build

WORKDIR /src

COPY shared ./shared
COPY stream ./stream

RUN go work init ./shared ./stream

WORKDIR /src/stream
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/stream ./cmd

FROM debian:bookworm-slim

RUN apt-get update && \
	apt-get install -q -y --no-install-recommends ffmpeg && \
	apt-get clean

COPY --from=build /out/stream /stream

EXPOSE 8083

ENTRYPOINT [ "/stream" ]
