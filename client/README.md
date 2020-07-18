# ponieproxy-client
A client of ponieproxy, which easily allows you to write your own proxy filters

## Structure
In the `internal` folder, lives the goproxy wrapper.

In the `client` folder, is the place where you can write filters.

In `cmd/main.go` you initiate ponieproxy and you add your request and response filters.

In `client/utils/utils.go` you can add your utility/helper functions.

In `client/filters/` you add your ponieproxy filters and view the default ones.

## Filters

## Usage
Clone this repository.

Run the client, by running:
`go run . -u URL_FILE -o OUTPUT_DIR`
