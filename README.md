# goComms repo
This repo is a communication framework over TCP.


# Stacks

## Bottom
## Top
## Message Breaker
## Compression

## WebSocket
Only a client websocket implemented. Server websocket must still be implemented

## TLS
Both client and server side can connect via TLS
### Todo
* Generic certificate is working, but need a generic Certificate service





# Example






# Todo
* Remove Stack definitions inside ./common
  * Each connection needs to define which stacks it will use and pass it on via optional params in the connection creation
    * This is currently implemented in the listen side, and used in an example







# Change Log


## Version 0.1.0
* Initial version

