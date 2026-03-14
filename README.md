## Web Cache

A simple project that simulates the role of a reverse proxy acting as a cache server.

To simplify the process, the client must include the desired URL in the **req_obj** JSON property of the request sent to the web cache. The web cache then responds to the client either with the cached response or with the response fetched from the requested URL, depending on whether a cache hit occurs.


## Sequence Diagram

```mermaid
sequenceDiagram
    Client->>+Web Cache: HTTP: GET
    participant Cache@{"type" : "database" }
    Web Cache->>+Cache: Has in Cache?
    alt has
        Web Cache->>Server: HTTP: Conditional-GET
        Server->>Web Cache: HTTP: 304 || 200
        alt cache valid
            Cache->>Client: HTTP: 200(Cache)
        else not valid
            Server->>Web Cache: (Revalidate Cache)
        end
    else does not have
        Web Cache->>Server: HTTP: GET
        Server->>Web Cache: HTTP: (Store in cache)
    end
    Web Cache->>Client: HTTP: 200
    
