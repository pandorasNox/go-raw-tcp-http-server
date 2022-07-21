# ...

## presentation idea
- disclaimer / expectations
  - ambitous - writing as much ourselfs as possible - so don't now how far we'll come
  - this is a more practical aproach to tcp rather than going into protocol details
  - its learning by using / doing it
  - to a certain degree I'm pretty much a bottom up guy
- why
  - orginally I stumbled upon rust tcp webserver (in the offical book)
    - https://doc.rust-lang.org/book/ch20-01-single-threaded.html
  - I though, ey I can do this golang
  - the other influence is
    - I always wanted to have like a man in the middle traffic sniffer
      - something like wireshark (but simpler?)
- so the jouney begins
- past
  - first time actively stumble upon tcp
    - k8s pod/deployment/service yamls
    - where you have to define the port protocol
      ``` yaml
      apiVersion: v1
      kind: Pod
      metadata:
        name: nginx
      spec:
        containers:
          image: nginx:1.23.1-alpine
          name: nginx
          ports:
          - containerPort: 80
            name: http
            protocol: TCP
      ```
      ``` yaml
      apiVersion: v1
      kind: Service
      metadata:
        name: my-service
      spec:
        selector:
          app: MyApp
        ports:
          - protocol: TCP
            port: 80
            targetPort: 9376
      ```
    - I assume its related to k8s routing capabilities
    - not in docker:
      ``` yaml
      # docker-compose.yaml
      services:
        web:
          build: .
          ports:
            - "8000:8000"
      ```
- tcp
  - what was that again?
    - Transmission Control Protocol
- some usefull infos in the course of this journey
  - often mentioned as part of the "Internet protocol suite" (IP/TCP/UDP)
  - it pairs with IP (internet protocol)
- simple tcp request response
  - text over tcp
  - or binary (will skipp that - as I couldn't find a cool example)
- ...
- http
  - curl making request forever
  - reading the http request
  - giving response
    - curl finishes
- tcp
  - http
    - server
      - go
      - rust
    - client
      - curl
  - custom tcp protocol (text based)
    - server
      - simple
        - go
        - node
        - php
    - client
      - go
- OSI
  - layer 4 vs 7
- ...
- tcp proxy sniffer?
  - src: https://github.com/jpillora/go-tcp-proxy
  - http
  - redis
    - `./tools redis-server`
    - proxy
      - `go run example/go/tcp/proxy/main.go -t localhost:6379 -l ":6479" -v`
    - `./tools redis-cli`
      - `redis-cli -h host.docker.internal -p 6379 ping`
      - after proxy
      ``` shell
      redis-cli -h host.docker.internal -p 6479 ping
      redis-cli -h host.docker.internal -p 6479 INFO
      redis-cli -h host.docker.internal -p 6479 set foo bar
      redis-cli -h host.docker.internal -p 6479 get foo bar
      redis-cli -h host.docker.internal -p 6479 get foo
      ```

example/go/tcp/proxy/main.go `
- todo:
  - telnet example

## resources
* https://www.youtube.com/watch?v=BHxmWTVFWxQ
* https://hackthedeveloper.com/golang-http-using-tcp-go-net-package/
* https://gobyexample.com/http-servers
* https://image4.io/en/blog/http-request-structure
* https://medium.com/@ahmet9417/golang-thread-pool-and-scheduler-434dd094715a
* https://gobyexample.com/signals
* https://stackoverflow.com/questions/19182367/go-http-server-connection-pool
* https://betterprogramming.pub/build-a-tcp-connection-pool-from-scratch-with-go-d7747023fe14
* https://itnext.io/forcefully-close-tcp-connections-in-golang-e5f5b1b14ce6?gi=f1878dec65f6#:~:text=Close()%20works,then%20close%20the%20TCP%20session. 
* https://riptutorial.com/node-js/example/22405/a-simple-tcp-server
* [Simple TCP echo server in Go
](https://www.youtube.com/watch?v=ijVQTdCGCqA)
* https://github.com/jpillora/go-tcp-proxy

## security
* https://www.youtube.com/watch?v=XiFkyR35v2Y&ab_channel=Computerphile

## resouce limits golang
* https://www.youtube.com/watch?v=nok0aYiGiYA&ab_channel=GopherAcademy
* https://groups.google.com/g/golang-nuts/c/rxxSBdkHOck
