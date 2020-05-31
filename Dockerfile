FROM golang:1.14 as backend
WORKDIR /build
COPY ./main.go .
COPY ./go.mod .
COPY ./go.sum .
RUN GOOS=linux GOARCH=amd64 go build -o main

FROM node:12.17-buster as frontend
WORKDIR /build
COPY ./interface/ .
RUN npm install
RUN npm run build

FROM debian:buster-slim
WORKDIR /app
COPY --from=backend /build/main .
COPY --from=frontend /build/dist ./dist
CMD /app/main
