FROM node:latest AS frontend

RUN mkdir /build
COPY ./ /build
WORKDIR /build
RUN cd frontend && yarn install && yarn build && mv build ../static

FROM golang:1.14 AS backend

RUN mkdir /build
COPY ./ /build
WORKDIR /build
RUN GOOS=linux go build -o app .

FROM debian:bullseye-slim


RUN apt-get update && \
    apt-get install -y ca-certificates procps && \
    apt-get autoclean -y && apt-get clean -y && \
    apt-get autoremove -y && rm -rf /var/lib/{apt,dpkg,cache,log} && \
    mkdir /app

COPY --from=frontend /build/static .
COPY --from=backend /build/app .

EXPOSE 8080
CMD ["./app"]
