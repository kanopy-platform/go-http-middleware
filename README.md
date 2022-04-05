# go-http-middleware

[![Build Status](https://drone.corp.mongodb.com/api/badges/kanopy-platform/go-http-middleware/status.svg)](https://drone.corp.mongodb.com/kanopy-platform/go-http-middleware)

`go-http-middleware` is a collection of generic HTTP middleware functions.

Refer to the [Contributing Guidelines](./CONTRIBUTING.md) to get started.

## Release Management

Upon merge to mainline use Github's interface to create a new release following standard [senver](https://semver.org/) standards.  e.g. `v1.2.3`. 

## Middleware

### Logging

* logrus - Provides standardized HTTP Access structured logging via the third-party logrus package.
