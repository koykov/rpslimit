# RPS limit

Collection of various RPS (rate) limiters.

Currently supports 6 limiters:

* [Fixed Window](fixed_window.go)
* [Sliding Log](sliding_log.go) and optimized [v2](sliding_log_v2.go) version
* [Token Bucket](token_bucket.go)
* [Leaky Bucket](leaky_bucket.go)
* [Realtime Counter](realtime_counter.go) based on [counter](https://github.com/koykov/counter) package

## Fixed Window

The main idea: atomic counter increases at every allow call and checks if counter less than limit. Counter resets every
second to zero. Fastest solution, but precision is so low.

## Sliding Log

The main idea: counter collects times of each allow call. At every call, obsoleted logs removes and new time appends.
Allow method returns `true` if log's length less than limit.

V2 version has the same idea but optimized to reduce allocations (use linked list and stack instead of slice).

## Token Bucket

The main idea: counter has a channel with limited capacity. Every allow call tries to extract a value from channel and
returns `true` on success. At every second/limit interval new value puts to the channel. Channel must be full at init.

## Leaky bucket

The main idea: opposite to Token Bucket solution - every allow call tries to put new value to the channel and return
`true` on success. At every second/limit interval one value extracts from the channel.

## Realtime counter

Based on [counter](https://github.com/koykov/counter) package. This solution just a wrapper over the counter. Every
allow call increases the counter and then check if sum is less than limit.
