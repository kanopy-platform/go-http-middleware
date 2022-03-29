# go-http-middleware

[![Build Status](https://drone.corp.mongodb.com/api/badges/kanopy-platform/go-http-middleware/status.svg)](https://drone.corp.mongodb.com/kanopy-platform/go-http-middleware)

The `go-http-middleware` is a collection of generic HTTP middleware functions.

## Design Guidelines

* Middleware MUST NOT rely on external application specific state
* Middleware MAY add information to the request context for downstream use.
* Generic Middleware MUST not interrupt the chain. It MUST always call the next handler.

Refer to [logrus](./logging/logrus.go) as an example of the middleware design.

## Logging

* logrus - Provides standardized HTTP Access structured logging via the third-party logrus package.
