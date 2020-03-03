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

FROM chromedp/headless-shell:80.0.3987.122

RUN apt-get update && \
    apt-get install -y dumb-init ca-certificates && \
    apt-get autoclean -y && apt-get clean -y && \
    apt-get autoremove -y && rm -rf /var/lib/{apt,dpkg,cache,log} && \
    mkdir /app

COPY --from=frontend /build/static ./static
COPY --from=backend /build/app .

EXPOSE 8080
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["./app"]
