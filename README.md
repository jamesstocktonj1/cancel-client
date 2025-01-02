# Wasmcloud HTTP Cancel Bug

This was discovered when running [`hey`](github.com/rakyll/hey) with a high number of concurrent requests. Occasionally when running with a low number of component instances it would stop responding to the requests. It would then only start responding when you run `wash stop component ...` or re-scale the number of instances.

On further investigation, I was able to reproduce the bug by cancelling the http request just after the response headers are being sent. This causes that current instance of the http component (`wasi:http/incoming-handler`) to hang indefinetly (or until the `max-execution-time` value is reached). 

I created this simple [Go script](https://github.com/jamesstocktonj1/cancel-client/blob/main/cmd/cancel-client/main.go) to try and reproduce this. The `sendCancelRequest()` will run a single cancelled request and stop one instance. The `countNumberInstances()` function will keep sending a cancelled request until the normal requests timeout, this gives an estimate of how many component instances are running within Wasmcloud.

I believe this is to do with the http provider dropping the http connection before consuming the response topic. Whilst the component is waiting for the provider to consume the topic, it does not know that the request was cancelled and the connection dropped. This may not be limited to just the http provider/component, if a component or provider calling a component drops the connection and does not consume the result, it may result in the same situation.

## Recreate
Start by deploying a fresh instance of [Wasmcloud](https://github.com/wasmCloud/wasmCloud):
```
$ wash up -d
```

Deploy the test application, a simple http-server + http-hello-world-rust example:
```
$ wash app deploy wadm.yaml
```

Run the cancel-client program and observe the following:
```
$ go run cmd/cancel-client/main.go 
2025/01/02 18:24:30 Cancel Request: Get "http://localhost:8080/": net/http: timeout awaiting response headers
2025/01/02 18:24:32 Normal Request: Get "http://localhost:8080/": context deadline exceeded
```

After running this then curl will hang when you try to send a request:
```
$ curl localhost:8080/
```