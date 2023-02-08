# go-libp2p-pubSub-stream

## Installation

Before you start you need to clone this repository.

Clone the go-libp2p-pubSub-stream repository:

```shell
git clone git@github.com:Parviz-Makhkamov/go-libp2p-pubSub-stream.git
```


## Usage

For start pubSub, in the root directory of the cloned repository:

From Alice side:
```shell
cd ./pubSub/alice-side
go run .  // started pubSub with default topic name and port number
or
go run . -topicName=myTopic, -listenAddrs=/ip4/127.0.0.1/tcp/PortNumber
```

From Bob side:
```shell
cd ./pubSub/bob-side
go run .  // started pubSub with default topic name and port number
or
go run . -topicName=myTopic, -listenAddrs=/ip4/127.0.0.1/tcp/PortNumber // with special topic name and port number
```

For start stream, in the root directory of the cloned repository:

From Alice side:
```shell
cd ./stream/alice-side
go run .  // started stream with default port number
or
go run . -l= PortNumber // started stream with special PortNumber
```

From Bob side:
```shell
cd ./stream/bob-side
go run . -l  PortNumber+AnyNumber -d /ip4/127.0.0.1/tcp/9007/p2p/QmRXhXyCApoYTi8esQv28awPV5MUyK7pW7AiPMGjcFNKq1
```
