package storage

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ActionRun struct {
	ID           bson.ObjectID  `bson:"_id,omitempty" json:"id"`
	Timestamp    time.Time      `bson:"timestamp" json:"timestamp"`
	Action       string         `bson:"action" json:"action"`
	Risk         string         `bson:"risk" json:"risk"`
	Mode         string         `bson:"mode" json:"mode"`
	Input        map[string]any `bson:"input,omitempty" json:"input,omitempty"`
	Outcome      string         `bson:"outcome" json:"outcome"`
	Summary      string         `bson:"summary,omitempty" json:"summary,omitempty"`
	ErrorMessage string         `bson:"error_message,omitempty" json:"errorMessage,omitempty"`
	DurationMS   int64          `bson:"duration_ms" json:"durationMs"`
	RequestID    string         `bson:"request_id,omitempty" json:"requestId,omitempty"`
}

func (m *MongoStorer) RecordActionRun(run ActionRun) (id string, err error) {
	ctx, cancel := m.newContext()
	defer cancel()

	if run.Timestamp.IsZero() {
		run.Timestamp = time.Now().UTC()
	}
	run.ID = bson.NewObjectID()

	if _, err = m.actionRuns.InsertOne(ctx, run); err != nil {
		return
	}
	id = run.ID.Hex()
	return
}

func (m *MongoStorer) ListActionRuns(limit int) (runs []ActionRun, err error) {
	ctx, cancel := m.newContext()
	defer cancel()

	if limit <= 0 || limit > 500 {
		limit = 100
	}

	cursor, err := m.actionRuns.Find(
		ctx,
		bson.M{},
		options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}}).SetLimit(int64(limit)),
	)
	if err != nil {
		return
	}
	defer cursor.Close(ctx)

	runs = []ActionRun{}
	err = cursor.All(ctx, &runs)
	return
}

func (m *MongoStorer) GetActionRun(id string) (run ActionRun, found bool, err error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		err = nil
		return
	}

	ctx, cancel := m.newContext()
	defer cancel()

	err = m.actionRuns.FindOne(ctx, bson.M{"_id": objectID}).Decode(&run)
	switch err {
	case mongo.ErrNoDocuments:
		err = nil
	case nil:
		found = true
	}
	return
}
