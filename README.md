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
    
