# Websocket
Details the commands used to interact with the websocket server (w/ examples).

## Commands
Each command must start with:
```
{
  "cmd": "command", //Command identifyer
  //... command specific requirements
}
```

## Auth
The auth command is used to initate the websocket client and identifty/authenticate accounts.
> This may also be used eventually as event encryption key instantiation 
```
{
  "cmd":"auth",
  "stream": 13122, //stream ID unique to the user/org (int32)
  "key": "iOMkjCb0", //Random string to identifiy as a svc account (string ~ 16 chars)
  "secret": "hVAF9UHsL68MgUBAGuc5lHodGv0yC9s5" //Random string acting as a svc account auth token (string ~ 32 chars)
}
```

## Subscribe
The subscribe commands instructs the websocket server to subscribe to a user stream channel (w/ NATS compatible channel wildcards)
> channels are prefixed in backend with `_USER.{stream_id}.` to maintain account seperation
```
{
  "cmd":"sub",
  "channel":"*.github.com" //The channel to subscribe to (w/ support for wildcards)
}
```

## Publish
The publish command instructs the websocket server to publish data into a channel
```
{
  "cmd":"pub",
  "channel": "envtsrc.github.com", //The channel to publish to
  "type":"repo_push", //The event type
  "data": { //The event data
    "branch":"master",
    "commit_id":"a94a8fe5ccb19ba61c4c0873d391e987982fbbd3",
    "author":"tcfw"
  },
  "typeVersion":"0.1" //Event version (optional)
}
```

## Replay
The replay command passes a request to start replaying events starting at a particular timestamp.
Replays are rebroadcasted to all subscribers through the related channel
```
{
  "cmd":"replay",
  "channel":"evntsrc.github.com", //Channel to replay from (no wildcards)
  "startTime": "2018-01-01 00:00:00 +1000" //Time to start replays from (flexible format using github.com/simplereach/timeutils)
  "filter": { //Filter the replay (optional)
    // TBC
  }
}
```
