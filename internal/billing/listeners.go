package billing

import (
	"context"
	"encoding/json"
	"log"

	"github.com/tcfw/evntsrc/internal/billing/events"
	userEvents "github.com/tcfw/evntsrc/internal/users/events"
	evntsrc_users "github.com/tcfw/evntsrc/internal/users/protos"
	"github.com/tcfw/evntsrc/internal/utils/sysevents"
)

func (s *Server) listenLoop() {
	newUsers, closeNewUsers, err := sysevents.ListenForBroadcast("billing", userEvents.BroadcastTypeCreated, "")
	defer closeNewUsers()
	if err != nil {
		return
	}

	for {
		select {
		case msg := <-newUsers:
			userEv := &userEvents.Event{}
			if err := json.Unmarshal(msg, userEv); err != nil {
				log.Printf("Failed to parse user event: %s", err)
				continue
			}

			go s.handleNewUser(userEv)
		}
	}
}

func (s *Server) handleNewUser(ev *userEvents.Event) {
	ctx := context.Background()
	userCli, err := newUserClient(ctx)
	if err != nil {
		log.Printf("Failed to open cli to users: %s", err)
		return
	}
	user, err := userCli.Find(ctx, &evntsrc_users.UserRequest{Status: evntsrc_users.UserRequest_PENDING, Query: &evntsrc_users.UserRequest_Id{Id: ev.UserID}})
	if err != nil {
		log.Printf("Failed to find user: %s", err)
		return
	}

	customer, err := s.createCustomer(user, "")
	if err != nil {
		log.Printf("Failed to create customer for new user: %s", err)
	} else {
		log.Printf("Created new customer (%s) for user (%s)", customer.ID, user.Id)
	}

	bEvent := &events.Event{Event: &sysevents.Event{Type: events.BroadcastTypeCreated}, UserID: ev.UserID, CustomerID: customer.ID}
	if err := sysevents.BroadcastEvent(context.Background(), bEvent); err != nil {
		log.Printf("Failed to broadcast creation event: %s", err)
	}
}
