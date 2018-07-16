# TODO's, improvements and notes

## "Layers"

- Shared Kernel: Common
- Domain: Endpoints
- Application: Datastore
- Framework: Main

## Config

All configuration is currently hardcoded into the source, which is put lightly - pretty bad.
Configuration should be loaded by envvar, flag, config file or by connecting to a service. This heavily depends on the use case, so it's been left out despite it being a critical part.

## TLS

The game client prototype is using unencrypted http communication, so TLS has been left out.

## Errors

Errors should include an error level (WARN, ERR, PANIC, FATAL) and a stack trace, to play into the logging better.

Error type should be separated into a base type, which isn't tied to a specific protocol, and then make the current version an extended HTTP version.

With the protocol specific extension, an optional payload can be added to give better error responses for the protocol.

## Handlers & Endpoints

Handlers and Endpoint currently do not conform to the Port and Adapters design pattern.
They should be reworked so that the endpoints can be tested in isolation of the choice of adapter (HTTP).

Add middleware to the endpoints, so it's possible to easily add in things such as metrics.

There's a lack of validation in the endpoints, because the current description is rather lackluster and touches nothing on the restrictions of the API, and it's rather vague on it's intent. And rather than assuming intent and setting up validation, I've chosen to "follow the instructions", or the lack thereof, and just say we can add the validation later.

## Datastore (Secondary Adapter)

Improved unit test cases.
Implement sql test case.

## Metrics and Status / Health

A Status/Health endpoint which checks the health of the service and all it's connections, which then can be polled by a monitoring service.

Add metric middleware to the endpoints and aggregate it in storage, where it can be analyzed and trigger warnings in case of larger deviations.

## Common Lib (Shared Kernel / Standard Lib / ... )

All written more or less specific for this use case and not generic enough to be used in a growing system. But it's a starting point to show the idea of a shared kernel, and simplifying the service itself.

Also, it feels a bit lonely without any tests in there.

[Back to README](README.MD)
