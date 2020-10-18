# HTTP Key Value Store

# Information

This is an implementation of an HTTP webserver acting as an in-memory key value store.
Supported Operations:
1. POST: /db/{key: <>, value: <value>} : Adds a new entry to the kv store, overwrites a previously stored key if any.
2. GET: / => Returns a json encoded {key: <>, value: <>} response if key is found in the map.
  
The server supports high throughput concurrent connections. Approximately 40k rps with wrk tool.
