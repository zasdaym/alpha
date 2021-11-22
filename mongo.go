package alpha

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "embed"

	"github.com/qiniu/qmgo"
	opts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoServer provides AlphaServer with MongoDB storage implementation.
type MongoServer struct {
	client *qmgo.QmgoClient
}

// NewMongoServer creates a new MongoServer.
func NewMongoServer(client *qmgo.QmgoClient) *MongoServer {
	return &MongoServer{
		client: client,
	}
}

// Handler returns HTTP handler for MongoServer.
func (ms *MongoServer) Handler() (*http.ServeMux, error) {
	mux := http.NewServeMux()
	indexHandler, err := ms.handleIndex()
	if err != nil {
		return nil, err
	}
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/increment", ms.handleIncrement())
	return mux, nil
}

//go:embed "views/index.html"
var indexHTML string

func (ms *MongoServer) handleIndex() (http.HandlerFunc, error) {
	tmpl, err := template.New("index").Parse(indexHTML)
	if err != nil {
		return nil, err
	}

	type attempt struct {
		ClientID string `bson:"client_id"`
		Count    int    `bson:"count"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var attempts []attempt
		if err := ms.client.Find(r.Context(), bson.D{{}}).All(&attempts); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		if err := tmpl.Execute(w, attempts); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}, nil
}

func (ms *MongoServer) handleIncrement() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientID := r.PostFormValue("client-id")
		if err := ms.increment(r.Context(), clientID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}
}

func (ms *MongoServer) increment(ctx context.Context, clientID string) error {
	if err := ms.client.UpdateOne(
		ctx,
		bson.M{"client_id": clientID},
		bson.M{
			"$inc": bson.M{"count": 1},
		},
		opts.UpdateOptions{
			UpdateOptions: options.Update().SetUpsert(true),
		},
	); err != nil {
		return fmt.Errorf("upsert attempt: %w", err)
	}
	return nil
}
