package ttlworker

import (
	"context"
	"io"
	"time"

	storer "github.com/tcfw/evntsrc/internal/storer/protos"
)

func (w *Worker) processBindings() error {
	return nil
}

func (w *Worker) processStream(stream int32) error {
	ctx := context.Background()

	streamEvents, err := w.storerCli.Query(ctx, &storer.QueryRequest{Stream: stream, Query: &storer.QueryRequest_Ttl{}})
	if err != nil {
		return err
	}

	for {
		event, err := streamEvents.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		//replay event

		//Extend TTL
		var cTTL *time.Time
		currentTTL, hasTTL := event.Metadata["ttl"]
		if hasTTL {
			parsed, _ := time.Parse(time.RFC3339, currentTTL)
			cTTL = &parsed
		}

		nTTL := time.Now().Add(5 * time.Minute)
		//TODO(tcfw) allow stream and event based TTL retry config

		_, err = w.storerCli.ExtendTTL(ctx, &storer.ExtendTTLRequest{Stream: stream, EventID: event.GetID(), CurrentTTL: cTTL, TTLTime: &nTTL})
		if err != nil {
			return err
		}
	}

	return nil
}
