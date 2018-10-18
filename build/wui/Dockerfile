FROM node:latest as builder

WORKDIR /web

COPY web/package.json package.json
COPY web/yarn.lock yarn.lock

RUN yarn 

COPY web /web
RUN yarn build

FROM nginx:latest

COPY --from=builder /web/dist /usr/share/nginx/html