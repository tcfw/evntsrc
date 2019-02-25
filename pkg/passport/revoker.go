package passport

import (
	"context"

	"github.com/globalsign/mgo/bson"
	pb "github.com/tcfw/evntsrc/pkg/passport/protos"
	"github.com/tcfw/evntsrc/pkg/utils/db"
	events "github.com/tcfw/evntsrc/pkg/utils/sysevents"
)

const dbName = "pssport"
const dbCollection = "revokes"

func isTokenRevoked(jti string) (bool, error) {
	dbConn, err := db.NewMongoDBSession()
	if err != nil {
		return true, err
	}
	defer dbConn.Close()

	collection := dbConn.DB(dbName).C(dbCollection)
	query := collection.Find(bson.M{"jti": jti})

	count, err := query.Count()
	if err != nil {
		return false, nil
	}

	return count != 0, nil
}

func revokeToken(claims map[string]interface{}, reason string) error {
	dbConn, err := db.NewMongoDBSession()
	if err != nil {
		return err
	}
	defer dbConn.Close()

	collection := dbConn.DB(dbName).C(dbCollection)

	revoke := &pb.Revoke{
		Id:     bson.NewObjectId().Hex(),
		Jti:    claims["jti"].(string),
		Reason: reason,
	}

	if err = collection.Insert(revoke); err != nil {
		return err
	}

	events.BroadcastEvent(context.Background(), &events.UserEvent{
		Event: &events.Event{
			Type:     "io.evntsrc.passport.revoked",
			Metadata: map[string]interface{}{"jti": claims["jti"]},
		},
		UserID: claims["sub"].(string),
	})

	return nil
}
