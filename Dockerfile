FROM golang:1.15-alpine AS build

ENV GO111MODULE=on

RUN apk update && \
    apk upgrade && \
    apk add --no-cache \
    git

RUN mkdir /app
COPY . /app
WORKDIR /app

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o consul-health-check

FROM mongo:3.6
COPY --from=build /app/ /usr/local/bin/

COPY consul-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/consul-entrypoint.sh
RUN chmod +x /usr/local/bin/consul-health-check

ENTRYPOINT ["consul-entrypoint.sh"]

EXPOSE 27017
CMD ["mongod"]