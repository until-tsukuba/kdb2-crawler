package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/oklog/ulid/v2"
)

type SubjectRef struct {
	CourseID string `json:"courseID"`
	Title    string `json:"title"`
}

type Subject struct {
	CourseID       string        `json:"courseID"`
	Title          string        `json:"title"`
	Credit         float32       `json:"credit"`
	Grade          int           `json:"grade"`
	Timetable      string        `json:"timeTable"`
	Books          []string      `json:"books"`
	ClassName      []string      `json:"className"`
	PlanPretopics  string        `json:"planPretopics"`
	Keywords       []string      `json:"keywords"`
	SeeAlsoSubject []*SubjectRef `json:"seeAlsoSubject"`
	Summary        string        `json:"summary"`
}

func main() {
	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %v", err)
		return
	}

	id := ulid.MustNew(ulid.Timestamp(time.Now()), ulid.DefaultEntropy())
	req := esapi.IndexRequest{
		Index:      "kdb2",
		Body:       os.Stdin,
		DocumentID: id.String(),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		var buf bytes.Buffer
		_, _ = buf.ReadFrom(res.Body)
		fmt.Fprintf(os.Stderr, "err resp: %v\n", buf.String())
		return
	}

	fmt.Printf("Document indexed successfully: %v\n", id.String())
}
