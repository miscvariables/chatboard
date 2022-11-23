# chatboard
This is a simple chat web app, you can use it to share your clipboard among devices.

## Why
Many times I wanta share simple information(text, url) between several devices: PC, Phone, TVbox, or share information with friends. I can't find a simple tool or web service to do that. So I made this.

## How
When you share text or url, you just type or copy&paste something to others, and other people can also do that. Yes, you are chatting! Clipboard sharing is just like chatting.

In short, This is a simple chat web app, you run the chat server yourself, or visit the server that I am running at: http://www.oddkits.com/chatboard

The code was manually forked from https://github.com/ezynda3/go-chat, I add some ideas, for example: Rooms

## How to use
It's very simple to use, visit http://www.oddkits.com/chatboard from your PC, set a Room and Username for yourself(Default Room: @room0, user: b), then open it again in your Phone, enter the SAME Room, then you can chat(share) between the PC and Phone. Everyone in the same Room can see same messages.


All dependencies are located in the vendor folder so Go 1.6+ is required.

## How to run it by yourself

Install bindata tools to generate bindata.so, this is optional, you only need to do this if you modified public/*
```
go get github.com/go-bindata/go-bindata/...
go get github.com/elazarl/go-bindata-assetfs/...
export PATH=$PATH:$HOME/go/bin
go-bindata-assetfs public
```

Build, install dependencies first
```
go get github.com/gorilla/websocket
go clean --cache
go build -o chatboard .
```

Run with default web port 80, or set the port yourself
```
./charboard
./charboard 8000
```

Then point your browser to http://localhost:8000
