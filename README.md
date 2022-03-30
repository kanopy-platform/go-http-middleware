# go-http-middleware

[![Build Status](https://drone.corp.mongodb.com/api/badges/kanopy-platform/go-http-middleware/status.svg)](https://drone.corp.mongodb.com/kanopy-platform/go-http-middleware)

The `go-http-middleware` is a collection of generic HTTP middleware functions.

## Design Guidelines

* Middleware MUST follow [loose coupling](https://en.wikipedia.org/wiki/Loose_coupling) principles
* Middleware MAY add information to the request context

Refer to [logrus](./logging/logrus.go) as an example of the middleware design.

## Logging

* logrus - Provides standardized HTTP Access structured logging via the third-party logrus package.
