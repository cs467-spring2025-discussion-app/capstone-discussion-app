The `auth` directory contains the code to build the user registration and authentication system for the platform.

## Build

Build like a typical Go module:

```bash
go build -o auth ./main.go
```

## Directory Structure

```plaintext
 auth
├──  docs
├──  internal
│   ├──  database
│   └──  testutils
├──  pkg
│   └──  logger
└──  scripts
```

`docs`: Contains documentation files related to the authentication system
`internal`: internal packages that are not meant to be used outside of the `auth` module

- `database`: code related to database interactions for the authentication system
- `testutils`: utility functions and types for testing the authentication system
- `pkg`: packages that are meant to be used by other modules
- `logger`: configuration and setup for logging
- `scripts`: utility scripts for local development and testing of the authentication system

## Dependencies

The `auth` module expects the following environment variables to be set:

- `DB`: The URL of the database to connect to
- `PORT`: The port to run the http server on
- `JWT_SECRET`: The secret key to sign the JWT tokens

See `example.env` or the `watch` command in `justfile` for sample environment variables.

Third party packages are defined in `go.mod` and `go.sum`.
