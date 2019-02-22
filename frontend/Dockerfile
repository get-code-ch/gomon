FROM node:10-alpine as builder
RUN apk add --no-cache --update git

RUN npm install -g @angular/cli

RUN git clone --progress --verbose https://github.com/get-code-ch/gomon-frontend.git
WORKDIR gomon-frontend

RUN rm -fv src/environments/*
COPY environments src/environments
RUN npm install --no-audit
RUN echo "npm install ended..."
RUN npm run build

FROM alpine:latest
MAINTAINER Claude Debieux <claude@get-code.ch>
RUN apk add --no-cache --update bash

COPY --from=builder /gomon-frontend/dist/gomon-frontend /client
CMD ["/bin/cp", "-vRf", "/client/.", "/gomon-frontend/"]
