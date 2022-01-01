# Occam Go Challenge project.

## Task

https://pastebin.com/DBnaYSDT

You are given an API interface of some service for streaming data, which is the price of a “Ticker” (eg. BTC price from an exchange). The price of BTC can be different on different exchanges, so our target is to build a “fair” price for BTC, combining the price from different sources. Let’s say there are up to 100 possible exchanges from where the price can be streamed.
You need to develop an algorithm which uses these data streams as input and as output providing an online “index price” in the form of minute bars, where the bar is a pair (timestamp, price). Output can be provided in any form (file, console, etc.). An example output if the service is working for ~2 minutes would look like:

```
Timestamp, IndexPrice
1577836800, 100.1
1577836860, 102
```

## Requirements

Data from the streams can come with delays, but strictly in increasing time order for each stream. Stream can return an error, in that case the channel is closed. Bars timestamps should be solid minute as shown in example and provided in on-line manner, price should be the most relevant to the bar timestamp. How to combine different prices into the index is up to you.
Code should be written using Go language, apart from that you are free to choose how and what to do. You can also write some mock streams in order to test your code.
The interface is artificial, so if you need to change something or to have additional assumptions - you are free to do this, but don’t forget to mention that. Your code will be reviewed but won’t be executed on our side.
We expect source code to be published on GitHub and shared with us followed by the readme file with the description of the solution written in English. There might be a technical call after the task completion where we can discuss the solution in detail, ask some questions etc.

## Technical solution
### Main considerations.

Firsts of all I have added the possibility to reconnect for the streams. 
If we have quite many sources-streams and one of them is disconnected, it is more optimal to reconnect the stream only, 
instead of restarting the application.

So I introduced the interface `RateSource` that is responsible for giving the last value from the stream plus the timestamp.
If the value is `nul` it means `undefined`. The stream is not working, and we do not take the value
into account. At the same time the stream is reconnected with retries and then the value will be available again.
The source should be closed for the cleanup of the resources.

### Changes of the streaming interface.

I have changed the `PriceStreamSubscriber` interface and the signature of `SubscribePriceStream` function. 
The output error channel was replaced by error. 
We receive the error if the stream cannot be connected and fails with error. If the error happens during the streaming 
(network), then the output channel is closed and the stream should be reconnected. Also, I added `Close`
function for final clean up of the resources (goroutines).

### Overall architecture

The main service - `IndexRateTicker` has the following components:
* `TimeTicker` for giving the time marks for the service. The implementation is `MinuteTimeTicker` giving timestamp on the beginning 
of every minute. This way we can change the time interval, providing different implementation
of the component.
* `RateSource` to have the last value of the rate form a stream or nil semantics for stream not available.
This component simplifies the process of getting data from stream to just local memory storage. And provides
the reconnection to the stream under the hood.
* `RateCalculator` for the component that calculate the mean price. This way we can have different strategies
for the mean value calculation for the rate. For example, with priorities for the sources-streams.

The components are created and assembled in `main` function, where we suppose to have the configuration.

Every component has a set of tests for the component. This makes the application being flexible
and well tested by unit tests.

### Conclusions.

Thus, we have configurable, flexible, well tested and understandable application. It works in concurrent mode,
closes on `syscall.SIGINT` and makes gentle shutdown of all its resources. The unit tests run with the --race flag, so we
have the protection against the race conditions and unexpected fails of the application in concurrency.

Also, the application is not provided with logging, tracing, metrics and the infrastructure things for the microservice
and the distributed systems. It does not have any real sources, just mocks for testing. And it has the console as the output.
So it does not have any real transport for the output. 
All that things are supposed to be outside the scope of the challenge exercise.


