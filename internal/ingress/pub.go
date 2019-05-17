package ingress

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/mux"
	"github.com/tcfw/evntsrc/internal/event"
)

//HandlePub publishes events to NATS
func HandlePub(w http.ResponseWriter, r *http.Request) {
	urlParts := mux.Vars(r)

	if _, ok := urlParts["stream"]; !ok {
		http.Error(w, "Unable to determine stream", http.StatusBadRequest)
		return
	}
	streamInt, _ := getStream(r)

	subject, ok := urlParts["subject"]
	if !ok {
		http.Error(w, "Unable to determine subject", http.StatusBadRequest)
		return
	}

	eventType, ok := urlParts["type"]
	if !ok {
		query := r.URL.Query()
		queryEventType := query.Get("type")
		if queryEventType == "" {
			http.Error(w, "Unable to determine event type", http.StatusBadRequest)
			return
		}
		eventType = queryEventType
	}

	pubEvent := &event.Event{}
	pubEvent.SetID()
	pubEvent.Stream = *streamInt
	pubEvent.Subject = subject
	pubEvent.Source = "ingress-webhook"
	pubEvent.CEVersion = "0.1"
	pubEvent.ContentType = r.Header.Get("Content-Type")
	pubEvent.Time = event.ZeroableTime{Time: time.Now()}
	pubEvent.Type = eventType
	pubEvent.TypeVersion = r.URL.Query().Get("version")
	pubEvent.Metadata = map[string]string{
		"source-ip": r.RemoteAddr,
	}

	userMetadata := r.URL.Query().Get("md")
	if userMetadata != "" {
		mdParts := strings.Split(userMetadata, " ")
		for key, part := range mdParts {
			pairs := strings.Split(part, ",")
			if len(pairs) == 2 {
				if _, ok := pubEvent.Metadata[pairs[1]]; ok {
					continue
				} else {
					pubEvent.Metadata[string(pairs[0])] = pairs[1]
				}
			} else {
				pubEvent.Metadata[string(key)] = pairs[0]
			}
		}
	}

	switch r.Method {
	case "GET":
		pubEvent.SetDataFromString(r.URL.Query().Get("data"))
		break
	case "POST":
		postData, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err == nil {
			pubEvent.Data = postData
		} else {
			log.Println(err.Error())
		}
		break
	}

	eventJSONBytes, _ := proto.Marshal(pubEvent.ToProtobuf())
	channel := fmt.Sprintf("_USER.%d.%s", pubEvent.Stream, pubEvent.Subject)

	natsConn.Publish(channel, eventJSONBytes)

	w.Write([]byte("Ok"))
}
