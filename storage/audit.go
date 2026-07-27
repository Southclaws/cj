package storage

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type AuditEvent struct {
	Timestamp  time.Time      `bson:"timestamp" json:"timestamp"`
	Action     string         `bson:"action" json:"action"`
	TargetType string         `bson:"target_type" json:"targetType"`
	TargetID   string         `bson:"target_id,omitempty" json:"targetId,omitempty"`
	Before     map[string]any `bson:"before,omitempty" json:"before,omitempty"`
	After      map[string]any `bson:"after,omitempty" json:"after,omitempty"`
	Outcome    string         `bson:"outcome" json:"outcome"`
	RequestID  string         `bson:"request_id,omitempty" json:"requestId,omitempty"`
}

func (m *MongoStorer) RecordAuditEvent(event AuditEvent) (err error) {
	ctx, cancel := m.newContext()
	defer cancel()

	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}

	_, err = m.audit.InsertOne(ctx, event)
	return
}

func (m *MongoStorer) ListAuditEvents(limit int) (events []AuditEvent, err error) {
	ctx, cancel := m.newContext()
	defer cancel()

	if limit <= 0 || limit > 500 {
		limit = 100
	}

	cursor, err := m.audit.Find(
		ctx,
		bson.M{},
		options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}}).SetLimit(int64(limit)),
	)
	if err != nil {
		return
	}
	defer cursor.Close(ctx)

	events = []AuditEvent{}
	err = cursor.All(ctx, &events)
	return
}
