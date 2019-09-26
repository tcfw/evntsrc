package users

import (
	"encoding/json"
	"log"

	"github.com/tcfw/evntsrc/internal/billing"
	billingEvents "github.com/tcfw/evntsrc/internal/billing/events"
	protos "github.com/tcfw/evntsrc/internal/users/protos"
	"github.com/tcfw/evntsrc/internal/utils/db"
	"github.com/tcfw/evntsrc/internal/utils/sysevents"
)

func (s *server) listenLoop() {
	billingCreated, closeBilling, err := sysevents.ListenForBroadcast("users", billingEvents.BroadcastTypeCreated, "")
	if err != nil {
		return
	}
	defer closeBilling()

	for {
		select {
		case billingCustomer := <-billingCreated:
			billingEv := &billingEvents.Event{}
			if err := json.Unmarshal(billingCustomer, billingEv); err != nil {
				log.Printf("failed to parse billing event: %s", err)
				continue
			}
			go s.handleBillingCreated(billingEv)
		}
	}
}

func (s *server) handleBillingCreated(bEv *billingEvents.Event) {
	dbConn, err := db.NewMongoDBSession()
	if err != nil {
		return
	}
	defer dbConn.Close()

	collection := dbConn.DB(dbName).C(dbCollection)
	query := collection.FindId(bEv.UserID)

	if c, err := query.Count(); c == 0 || err != nil {
		log.Printf("Failed to find user %s: %v", bEv.UserID, err)
		return
	}

	user := &protos.User{}
	if err = query.One(user); err != nil {
		log.Printf("Failed to build user %s: %v", bEv.UserID, err)
		return
	}

	user.Metadata[billing.UserMetadataCustomerID] = []byte(bEv.CustomerID)
	if err := collection.UpdateId(user.Id, user); err != nil {
		log.Printf("failed to update user with customer ID: %s", err)
	}
}
