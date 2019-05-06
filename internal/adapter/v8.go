package adapter

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/augustoroman/v8"
	"github.com/tcfw/evntsrc/internal/adapter/protos"
	"github.com/tcfw/evntsrc/internal/event/protos"
)

//RunV8Adapter takes in an adapter and executes the code in a V8 VM
func RunV8Adapter(s *Server, adapter *evntsrc_adapter.Adapter, srcEvent *evntsrc_event.Event) (*evntsrc_event.Event, []string, error) {

	if adapter.GetEngine() != evntsrc_adapter.Adapter_V8 {
		return nil, nil, fmt.Errorf("wrong adapter type passed")
	}

	ctx := <-s.v8Pool

	jsEvent, err := ctx.Create(srcEvent)
	if err != nil {
		return nil, nil, err
	}

	err = ctx.Global().Set("event", jsEvent)
	if err != nil {
		return nil, nil, err
	}

	log := []string{}
	done := make(chan error, 1)

	ctx.Global().Set("DataAsJsonString", ctx.Bind("DataAsJsonString", func(in v8.CallbackArgs) (*v8.Value, error) {
		return ctx.Create(string(srcEvent.Data))
	}))

	ctx.Global().Set("log", ctx.Bind("log", func(in v8.CallbackArgs) (*v8.Value, error) {
		log = append(log, in.Arg(0).String())
		return nil, nil
	}))

	ctx.Global().Set("cancel", ctx.Bind("cancel", func(in v8.CallbackArgs) (*v8.Value, error) {
		done <- fmt.Errorf("Event Cancelled")
		return nil, nil
	}))

	go func() {
		ret, err := ctx.Eval(string(adapter.GetCode()), "adapt.js")
		if err != nil {
			done <- err
		} else if ret.IsKind(v8.KindBoolean) && ret.Bool() == false {
			done <- fmt.Errorf("Adapt terminated falsy")
		} else {
			done <- nil
		}
	}()

	select {
	case <-time.After(10 * time.Second):
		ctx.Terminate()
		return nil, log, fmt.Errorf("Execution timeout")
	case err := <-done:
		if err != nil {
			ctx.Terminate()
			return nil, log, err
		}
	}

	retEvent, err := ctx.Global().Get("event")
	if err != nil {
		return nil, log, err
	}

	retEventBytes, err := retEvent.MarshalJSON()
	if err != nil {
		return nil, log, nil
	}

	adaptedEvent := &evntsrc_event.Event{}
	err = json.Unmarshal(retEventBytes, adaptedEvent)
	if err != nil {
		return nil, log, nil
	}

	ctx.Terminate()

	adaptedEvent = recoverImmutableEventData(srcEvent, adaptedEvent)

	return adaptedEvent, log, nil
}

func recoverImmutableEventData(srcEvent *evntsrc_event.Event, adaptedEvent *evntsrc_event.Event) *evntsrc_event.Event {

	adaptedEvent.ID = srcEvent.ID
	adaptedEvent.Stream = srcEvent.Stream
	adaptedEvent.CEVersion = srcEvent.CEVersion
	adaptedEvent.Time = srcEvent.Time

	for metaK, metaV := range srcEvent.Metadata {
		adaptedEvent.Metadata[metaK] = metaV
	}

	return adaptedEvent
}
