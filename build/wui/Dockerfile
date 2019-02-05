FROM node:latest as builder

WORKDIR /web

COPY web/package.json package.json
COPY web/yarn.lock yarn.lock

RUN yarn 

COPY web /web
RUN yarn build --modern

FROM nginx:latest

#Enable GZIP
COPY build/wui/compression.conf /etc/nginx/conf.d/

COPY --from=builder /web/dist /usr/share/nginx/html