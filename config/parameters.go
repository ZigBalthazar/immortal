package config

import (
	"context"
	"errors"

	"github.com/dezh-tech/immortal"
	"github.com/dezh-tech/immortal/database"
	"github.com/dezh-tech/immortal/handler"
	"github.com/dezh-tech/immortal/server/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Retention struct {
	Time  int         `bson:"time,omitempty"  json:"time,omitempty"`
	Count int         `bson:"count,omitempty" json:"count,omitempty"`
	Kinds interface{} `bson:"kinds,omitempty" json:"kinds,omitempty"`
}

type Subscription struct {
	Amount int    `bson:"amount" json:"amount"`
	Unit   string `bson:"unit"   json:"unit"`
	Period int    `bson:"period" json:"period"`
}

type Admission struct {
	Amount int    `bson:"amount" json:"amount"`
	Unit   string `bson:"unit"   json:"unit"`
}

type Publication struct {
	Kinds  []int  `bson:"kinds"  json:"kinds"`
	Amount int    `bson:"amount" json:"amount"`
	Unit   string `bson:"unit"   json:"unit"`
}

type Fees struct {
	Subscription []Subscription `bson:"subscription,omitempty" json:"subscription,omitempty"`
	Publication  []Publication  `bson:"publication,omitempty"  json:"publication,omitempty"`
	Admission    []Admission    `bson:"admission,omitempty"    json:"admission,omitempty"`
}

type Parameters struct {
	Handler         *handler.Config   `bson:"handler"                   json:"handler"`
	WebsocketServer *websocket.Config `bson:"server"                    json:"server"`
	Retention       *Retention        `bson:"retention,omitempty"       json:"retention,omitempty"`
	Fees            *Fees             `bson:"fees,omitempty"            json:"fees,omitempty"`
	Name            string            `bson:"name"                      json:"name"`
	Description     string            `bson:"description"               json:"description"`
	Pubkey          string            `bson:"pubkey"                    json:"pubkey"`
	Contact         string            `bson:"contact"                   json:"contact"`
	Software        string            `bson:"software"                  json:"software"`
	SupportedNips   []int             `bson:"supported_nips"            json:"supported_nips"`
	Version         string            `bson:"version"                   json:"version"`
	RelayCountries  []string          `bson:"relay_countries,omitempty" json:"relay_countries,omitempty"`
	LanguageTags    []string          `bson:"language_tags,omitempty"   json:"language_tags,omitempty"`
	Tags            []string          `bson:"tags,omitempty"            json:"tags,omitempty"`
	PostingPolicy   string            `bson:"posting_policy,omitempty"  json:"posting_policy,omitempty"`
	PaymentsURL     string            `bson:"payments_url,omitempty"    json:"payments_url,omitempty"`
	Icon            string            `bson:"icon,omitempty"            json:"icon,omitempty"`
}

func (c *Config) LoadParameters(db *database.Database) error {
	coll := db.Client.Database(db.DBName).Collection("config_parameters")

	ctx, cancel := context.WithTimeout(context.Background(), db.QueryTimeout)
	defer cancel()

	filter := bson.M{}

	var result *Parameters
	err := coll.FindOne(ctx, filter).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		// insert default parameters
		// for more info : https://github.com/nostr-protocol/nips/blob/master/11.md
		newDocument := &Parameters{
			Name:           "immortal",
			Description:    "a nostr relay designed for scale.",
			Pubkey:         "aca682c51c44c9046461de0cb34bcc6338d5562cdf9062aee9c3ca5a4ca0ab3c",
			Software:       "https://github.com/dezh-tech/immortal",
			SupportedNips:  []int{1, 11},
			Version:        immortal.StringVersion(),
			RelayCountries: []string{"*"},
			LanguageTags:   []string{"*"},
			Tags:           []string{},
			PostingPolicy:  "",
			PaymentsURL:    "",
			Icon:           "https://raw.githubusercontent.com/dezh-tech/immortal/refs/heads/main/assets/images/immortal.png",
			WebsocketServer: &websocket.Config{
				Limitation: &websocket.Limitation{
					MaxMessageLength:    8192,  // Maximum length of a single message (in bytes or characters)
					MaxSubscriptions:    20,    // Maximum number of concurrent subscriptions a client can create
					MaxFilters:          20,    // Maximum number of filters a client can apply in a subscription
					MaxSubidLength:      256,   // Maximum length of a subscription identifier
					MinPowDifficulty:    0,     // Minimum proof-of-work difficulty for publishing events
					AuthRequired:        false, // Whether authentication is required for writes
					PaymentRequired:     false, // Whether payment is required to interact with the relay
					RestrictedWrites:    false, // Whether writes are restricted to authenticated or paying users
					MaxEventTags:        1000,  // Maximum number of tags allowed in a single event
					MaxContentLength:    4096,  // Maximum content length of an event (in bytes)
					CreatedAtLowerLimit: 0,     // Earliest timestamp allowed for event creation
					CreatedAtUpperLimit: 0,     // Latest timestamp allowed for event creation (0 for no limit)
				},
			},
			Handler: &handler.Config{
				InitialQueryDefaultLimit: 100,
				Limitation: handler.Limitation{
					MaxLimit: 1000, // Maximum number of events returned in a query
				},
			},
			Retention: &Retention{},
			Fees: &Fees{
				Subscription: []Subscription{},
				Publication:  []Publication{},
				Admission:    []Admission{},
			},
			Contact: "",
		}

		insertErr := c.SetParameters(db, newDocument)
		if insertErr != nil {
			return insertErr
		}

		return nil
	} else if err != nil {
		return err
	}

	_, updateErr := coll.UpdateOne(ctx, filter, bson.D{{"$set", bson.D{{"version", immortal.StringVersion()}}}}) //nolint

	if updateErr != nil {
		return updateErr
	}

	c.Parameters = result
	c.WebsocketServer.Limitation = result.WebsocketServer.Limitation

	return nil
}

func (c *Config) SetParameters(db *database.Database, params *Parameters) error {
	coll := db.Client.Database(db.DBName).Collection("config_parameters")

	ctx, cancel := context.WithTimeout(context.Background(), db.QueryTimeout)
	defer cancel()
	_, insertErr := coll.InsertOne(ctx, params)
	if insertErr != nil {
		return insertErr
	}

	c.Parameters = params
	c.WebsocketServer.Limitation = params.WebsocketServer.Limitation

	return nil
}
