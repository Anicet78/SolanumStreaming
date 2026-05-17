FROM golang:1.25.8 AS build

WORKDIR /src

COPY shared ./shared
COPY movies ./movies

RUN go work init ./shared ./movies

WORKDIR /src/movies
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/movies ./cmd

FROM gcr.io/distroless/base-debian12:nonroot

COPY --from=build /out/movies /movies

EXPOSE 8082

USER nonroot:nonroot
ENTRYPOINT [ "/movies" ]
