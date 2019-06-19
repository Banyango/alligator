# Alligator

Alligator is a configurable reverse proxy in golang.

- Toml config files
- In memory Caching middleware

alligator.toml.sample contains an example of how to set up the proxies
```toml
[[Proxy]]
Path = "/whyhellothere"
Host = "localhost:5000"
Scheme= "http"
	[[Proxy.Rules]]
		Type = "host"
		Pattern = ["[localhost:3000]"]
	[[Proxy.Rules]]
		Type = "path"
		Pattern = [".hello/."]
```
This example redirects any host that matches the regex localhost:3000 and any path that contains hello/ to localhost:5000/whyhellothere

Things on my todo
- The proxy rules are effectively all AND which might be restrictive. There would be some rework required to handle OR, XOR, etc
- The https schemes are not fully tested, so I can't guarantee that it works.
- Using an in-memory cache that I did means you can't share the cache if you horz scale the proxy server.  
- In general there could be a lot more testing but I only did this in a day. 
- There are only 3 types of matchers there's probably more to match on.
- Didn't have time to setup CircleCI

###### To build the docker image
make v=1 build-docker

###### To run the docker image
`docker run -v $(pwd)/.dist/alligator.toml:/alligator.toml alligator:1`
