FROM golang:1.19.4-alpine3.17 as build

WORKDIR /app
COPY ./* /app/
RUN go build -o http2mqtt

FROM alpine:3.17.0 as runtime

WORKDIR /app
COPY --from=build /app/http2mqtt /app/

CMD ["/app/http2mqtt"]