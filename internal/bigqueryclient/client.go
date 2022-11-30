// Package bigqueryclient is responsible for interacting with the BigQuery API.
package bigqueryclient

import (
	"context"

	"cloud.google.com/go/bigquery"
)

type BigQueryClient interface {
	CreateTable(table string, schema interface{}) error
	Insert(table string, rows interface{}) error
}

type bigqueryClient struct {
	client  *bigquery.Client
	dataset string
	context context.Context
}

func (bq *bigqueryClient) CreateTable(table string, schema interface{}) error {
	s, err := bigquery.InferSchema(schema)
	if err != nil {
		return err
	}

	t := bq.client.Dataset(bq.dataset).Table(table)

	if err := t.Create(bq.context, &bigquery.TableMetadata{Schema: s}); err != nil {
		return err
	}

	return nil
}

func (bq *bigqueryClient) Insert(table string, rows interface{}) error {
	inserter := bq.client.Dataset(bq.dataset).Table(table).Inserter()
	if err := inserter.Put(bq.context, rows); err != nil {
		return err
	}

	return nil
}

func NewBigQueryClient(ctx context.Context, projectID string, dataset string) (BigQueryClient, error) {
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &bigqueryClient{client: client, dataset: dataset, context: ctx}, nil
}
