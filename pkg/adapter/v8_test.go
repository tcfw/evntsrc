package adapter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tcfw/evntsrc/pkg/adapter/protos"
	"github.com/tcfw/evntsrc/pkg/event/protos"
)

func TestRunAdapterV8BitbucketExample(t *testing.T) {
	adapter := &evntsrc_adapter.Adapter{
		ID:     "",
		Engine: evntsrc_adapter.Adapter_V8,
		Code: []byte(`
		let d = JSON.parse(DataAsJsonString());
		let subject = d.push.changes[0].new.name;
		event.eventSubject = subject;
    event.eventType = "push";
		`),
	}

	eventTime := time.Now()

	event := &evntsrc_event.Event{
		ID:          "test-id-00-00-0000000",
		Stream:      1,
		Time:        &eventTime,
		Type:        "test",
		CEVersion:   "0.1",
		Source:      "ingress-adapters",
		Subject:     "TEST-001",
		ContentType: "application/json",
		Data: []byte(`
{
  "push": {
    "changes": [
      {
        "forced": false,
        "old": {
          "type": "branch",
          "name": "TEST-1622"
        },
        "truncated": false,
        "commits": [],
        "created": false,
        "closed": false,
        "new": {
          "default_merge_strategy": "merge_commit",
          "merge_strategies": [
            "merge_commit",
            "squash",
            "fast_forward"
          ],
          "type": "branch",
          "name": "TEST-1622"
        }
      }
    ]
  },
  "repository": {
    "scm": "git",
    "name": "Some Repo",
    "project": {},
    "full_name": "some-company/repo",
    "owner": {},
    "type": "repository",
    "is_private": true,
    "uuid": "{00000000-0000-0000-0000-000000000000}"
  },
  "actor": {}
}
		`),
	}

	retEvent, _, err := RunV8Adapter(adapter, event)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "TEST-1622", retEvent.Subject)
}

func TestRunAdapterV8Cancel(t *testing.T) {
	adapter := &evntsrc_adapter.Adapter{
		ID:     "",
		Engine: evntsrc_adapter.Adapter_V8,
		Code:   []byte(`cancel();`),
	}

	eventTime := time.Now()

	event := &evntsrc_event.Event{
		ID:          "test-id-00-00-0000000",
		Stream:      1,
		Time:        &eventTime,
		Type:        "test",
		CEVersion:   "0.1",
		Source:      "ingress-adapters",
		Subject:     "TEST-001",
		ContentType: "application/json",
		Data:        []byte(""),
	}

	_, _, err := RunV8Adapter(adapter, event)

	assert.Error(t, err, "An error should have been thrown")
	assert.EqualError(t, err, "Event Cancelled")
}

func TestRunAdapterV8Log(t *testing.T) {
	adapter := &evntsrc_adapter.Adapter{
		Engine: evntsrc_adapter.Adapter_V8,
		Code:   []byte(`log("hi");`),
	}

	event := &evntsrc_event.Event{}

	_, log, err := RunV8Adapter(adapter, event)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, log)
	assert.Equal(t, "hi", log[0])
}

func TestRunAdapterV8ImmutableOverride(t *testing.T) {
	adapter := &evntsrc_adapter.Adapter{
		Engine: evntsrc_adapter.Adapter_V8,
		Code: []byte(`
    event.eventID = "foobar"
    event.eventStream = -1000;
    event.extensions.original = "false"; //will be reverted back
    event.extensions.somethingelse = "something";
    `),
	}

	eventTime := time.Now()

	event := &evntsrc_event.Event{
		ID:          "test-id-00-00-0000000",
		Stream:      1,
		Time:        &eventTime,
		Type:        "test",
		CEVersion:   "0.1",
		Source:      "ingress-adapters",
		Subject:     "TEST-001",
		ContentType: "application/json",
		Metadata: map[string]string{
			"original": "true",
		},
		Data: []byte(""),
	}

	retEvent, _, err := RunV8Adapter(adapter, event)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, retEvent.ID, "test-id-00-00-0000000")
	assert.Equal(t, retEvent.Stream, int32(1))
	assert.Equal(t, "true", retEvent.Metadata["original"])
}

func TestRunAdapterV8SyntaxError(t *testing.T) {
	adapter := &evntsrc_adapter.Adapter{
		Engine: evntsrc_adapter.Adapter_V8,
		Code: []byte(`
    cannot be processed;
    `),
	}

	event := &evntsrc_event.Event{}

	_, _, err := RunV8Adapter(adapter, event)

	assert.Error(t, err, "An exception should have occured")
	assert.Equal(t, "Uncaught", string(err.Error()[0:8]))
}

func TestRunAdapterV8Timeout(t *testing.T) {
	adapter := &evntsrc_adapter.Adapter{
		Engine: evntsrc_adapter.Adapter_V8,
		Code: []byte(`
function wait(ms){
   var start = new Date().getTime();
   var end = start;
   while(end < start + ms) {
     end = new Date().getTime();
  }
}

//wait 20 seconds
wait(20000);
    `),
	}

	event := &evntsrc_event.Event{}

	_, _, err := RunV8Adapter(adapter, event)

	assert.EqualError(t, err, "Execution timeout")
}

func TestRunAdapterV8AttemptLoadExternalJSviaDom(t *testing.T) {
	adapter := &evntsrc_adapter.Adapter{
		Engine: evntsrc_adapter.Adapter_V8,
		Code: []byte(`
var script = document.createElement('script');
script.src = "https://js.stripe.com/v3/";
document.head.appendChild(script);
    `),
	}

	event := &evntsrc_event.Event{}

	_, _, err := RunV8Adapter(adapter, event)

	assert.Error(t, err, "An exception should have occured")
}

func TestRunAdapterV8AttemptXHR(t *testing.T) {
	adapter := &evntsrc_adapter.Adapter{
		Engine: evntsrc_adapter.Adapter_V8,
		Code: []byte(`
var xhttp = new XMLHttpRequest()
xhttp.open("GET", "https://google.com.au", true)
xhttp.send();
    `),
	}

	event := &evntsrc_event.Event{}

	_, _, err := RunV8Adapter(adapter, event)

	assert.Error(t, err, "An exception should have occured")
}
