# consul-mongo

This dockerfile wraps the official mongo:3.6 docker image.
It adds an automatic service registration for [consul](https://www.consul.io) as well as a simple health check.
So everything from the [official mongo docker image](https://hub.docker.com/_/mongo) works exactly the same.

To use it you just have to pass the environment variable `CONSUL_HOST` which you have to set to the host where consul can be found:
```
CONSUL_HOST=consul-service:8500
```

The health check is running on a simple http webserver with the endpoint
```
GET :8101/healthcheck
```
If it returns the status code 200 everything is ok, else not.