# Stockplay
pure golang api server using alphavantage stock api.

## How to run
- make sure to install docker and docker-compose on your machine
- set `ALPHAVANTAGE_KEY` on file `docker-compose.yaml` to your `alphavantage`'s api key
- run `docker-compose up`
- the stock api server will be available at port `:8080` on your host machine
- for example `curl --request GET --url 'http://localhost:8080/?symbol=IBM'` will fetch `IBM` stock


## How to run the tests
- just run `go test ./..` in the project folder


## Summary
- Some stock market words still foreign to me.
- Connecting the 2 services feels repetitive and error prone, if it's real world project I would like to use
more defined format like `protobuf` (with or without grpc).
- I like the idea of monorepo in this kind of project, the authors of the service can maintain client
libraries and be responsible of every change in the service. If there any major breaking changes in the service
we can catch early if dependant services break.

