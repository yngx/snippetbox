# The http.Handler Interface

A handler is any object that satifies the http.Handler interface:

```
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

This means a handler object must have a ServeHTTP method with the following signature:

```
ServeHTTP(http.ResponseWriter, *http.Request)
```

Instead of creating an object to implement a ServeHTTP method. We can write a handler as a normal function.

```
func example(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("This is my example page"))
}
```

Then, to make it a handler, we transform it via tha `http.HandlerFunct()` adapter.

```
mux := http.NewServeMux()
mux.Handle("/", http.HandlerFunc(example))
```

The http.HandlerFunc() adapter works by automatically adding a ServeHTTP() method to the home function. When executed, this ServeHTTP() method then simply calls the content of the original home function. It’s a roundabout but convenient way of coercing a normal function into satisfying the http.Handler interface.

#### Another note:
The http.ListenAndServe() function takes a http.Handler object as the second parameter but we’ve been passing in a servemux.

```
func ListenAndServe(addr string, handler Handler) error
```

We were able to do this because the servemux also has a ServeHTTP() method, meaning that it too satisfies the http.Handler interface.

Just think of the servemux as just being a special kind of handler, which instead of providing a response itself passes the request on to a second handler. 


#### Another another note:

All incoming HTTP requests are served in their own goroutine. This means it’s very likely that the code in or called by your handlers will be running concurrently. Be weary of race conditions.