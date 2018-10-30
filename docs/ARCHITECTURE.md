

# Components
## Ingress
- http -> adapters
- ws
- grpc

## Store
- mongodb

## Cluster Distribution
- NATS

## Adapters
- JS V8 engine
- Python
- Lua
- Built in adapters:
  - bitbucket

## Egress
- http -> webhooks
- ws
- grpc streaming

## Secondary Features
- Persistence
- Replay

# Event structure
[CE Version 0.1](https://github.com/cloudevents/spec/blob/v0.1/spec.md) + Tenanted
```
type Event struct {
	Stream		int32 `json: omit`
	ID		string `json: "eventId"`
	Time		time.Time `json: "eventTime"`
	Type		string `json: "eventType"` //(rev-URI)
	TypeVersion	string `json: "eventTypeVersion"` //(1.0)
	CEVersion	string `json: "cloudEventVersion"` //(0.1)
	Source		string `json: "eventSource"` //(URI?) 
	Subject		string `json: "eventSubject"` //(optional ~ e.g. User ID / Task Number / Chat group etc)
	Acknowledged	time.Time `json: "eventAcknowledged"
	Metadata	map[string]interface{} `json: "extensions"` //(optional)
	ContentType	string `json: "contentType"` //(MIME~optional)
	Data		[]bytes `json: "data"` //(optional) 
}
```

# Data flow
1. Ingress
2. Process
	- Distribute -> NATS
  	- Store -> Mongo
3. Egress <- NATS