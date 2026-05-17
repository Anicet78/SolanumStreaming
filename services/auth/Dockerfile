FROM golang:1.25.8 AS build

WORKDIR /src

COPY auth ./auth
COPY shared ./shared

RUN go work init ./shared ./auth

WORKDIR /src/auth
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/auth ./cmd

FROM gcr.io/distroless/base-debian12:nonroot

COPY --from=build /out/auth /auth

EXPOSE 8081

USER nonroot:nonroot
ENTRYPOINT [ "/auth" ]
