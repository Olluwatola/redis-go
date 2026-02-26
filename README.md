# Redis Clone in Go

A lightweight Redis clone implementation in Go, capable of handling basic commands like `PING`, `SET`, `GET`, and `ECHO`.

## Features

- RESP (Redis Serialization Protocol) parsing
- Thread-safe in-memory key-value storage
- Key expiration support (EX/PX flags)
- Concurrent connection handling


## Running

```sh
./your_program.sh
```

The server listens on port 6379.

## Supported Commands

- `PING` - Returns PONG
- `ECHO <message>` - Returns the message
- `SET <key> <value> [EX seconds | PX milliseconds]` - Set a key with optional expiration
- `GET <key>` - Get a key's value

## SET Examples

```
SET key value              → No expiration
SET key value EX 10        → Expire in 10 seconds
SET key value PX 5000      → Expire in 5000 milliseconds
SET key value EXAT 1705315200  → Expire at Unix timestamp (seconds)
SET key value PXAT 1705315200000  → Expire at Unix timestamp (milliseconds)
```
