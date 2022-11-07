FROM golang:1.17-alpine as build
WORKDIR /app
RUN apk update && apk add --no-cache tzdata ca-certificates gcc musl-dev libc-dev git
COPY . .
RUN go mod download && mkdir /builds
RUN go build -ldflags '-extldflags "-fno-PIC -static"' -buildmode pie -tags 'osusergo netgo static_build' -o /builds/api

FROM alpine:latest as production
WORKDIR /app
COPY --from=build /usr/share/zoneinfo/ /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /builds/* ./
RUN chmod 755 ./*
CMD ["/app/api"]