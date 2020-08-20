# requests
A package utility that lets me make various requests in Go

```
Disclaimer: This is just a collection of abstractions that I use for my own personal projects
and is not meant to be a "serious" Go Package. The code documented here is for my own
benefit and learning.
```

### The Requestor Interface
```go
type Requestor interface {
    Request()
}
```
A `Requestor` is any object type that implements a Request() method to make HTTP requests.
Types that satisfy this interface are allowed to be passed into RequestPools

```go
type TwitterUser struct {
    username string
}

type(tu TwitterUser) Request() {
    // code that makes http request
}
```
The TwitterUser type can now be used as a `Requestor`


### Request Pool
```go
type RequestPool struct {
	wg         *sync.WaitGroup
	requestors []Requestor
	channel    chan Requestor
}
```
A `RequestPool` represents a pool of concurrent HTTP requests. When given an array of Requestors,
the RequestPool will spawn a worker thread for each Requestor present. The Requestor will then
make its request while inside the worker thread.

### Usage
```go
import pool "github.com/alexbotello/requests/pool"

func main() {
    requestors := []pool.Requestor{&TwitterUser{username: "BarackObama"}, &TwitterUser{username: "katyperry"}}
    // instantiate a new RequestPool
    rp := pool.NewRequestPool(requestors)
    // start the RequestPool -- goroutine will block here until all worker threads complete 
    rp.Start()
}
```
### Examples
[Twitter: Getting tweets](_examples/twitter.go)

RequestPool's implementation makes use of a common concurrency pattern called Fan-Out/Workers.
A visualization of the pattern looks like this:
![](animation.gif)

