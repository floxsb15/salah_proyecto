FROM golang:1.25.12-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/api .

FROM alpine:3.22

RUN apk add --no-cache ca-certificates && addgroup -S app && adduser -S -G app app
WORKDIR /app
RUN mkdir -p internal/images/usuarios internal/images/vehiculos internal/images/garantias && chown -R app:app /app
COPY --from=build --chown=app:app /out/api /app/api

USER app
EXPOSE 5000
CMD ["/app/api"]
