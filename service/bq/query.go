package bq

import (
	"context"
	"fmt"
	"bqtail/base"
	"bqtail/task"
	"google.golang.org/api/bigquery/v2"
)

func (s *service) Query(ctx context.Context, request *QueryRequest) (*bigquery.Job, error) {
	if err := request.Init(s.projectID); err != nil {
		return nil, err
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}
	job := &bigquery.Job{
		Configuration: &bigquery.JobConfiguration{
			Query: &bigquery.JobConfigurationQuery{
				Query:            request.SQL,
				UseLegacySql:     &request.UseLegacy,
				DestinationTable: request.destTable,
			},
		},
	}
	if request.Append {
		job.Configuration.Query.WriteDisposition = "WRITE_APPEND"
	}
	if request.UseLegacy {
		job.Configuration.Query.AllowLargeResults = true
	}

	if request.DatasetID != "" {
		job.Configuration.Query.DefaultDataset = &bigquery.DatasetReference{
			DatasetId: request.DatasetID,
			ProjectId: request.ProjectID,
		}
	}
	job.JobReference = request.jobReference()
	return s.Post(ctx, request.ProjectID, job, &request.Actions)
}

type QueryRequest struct {
	DatasetID string
	SQL       string
	UseLegacy bool
	Append    bool
	Dest      string
	destTable *bigquery.TableReference
	Request
}

func (r *QueryRequest) Init(projectID string) (err error) {
	if r.ProjectID != "" {
		projectID = r.ProjectID
	} else {
		r.ProjectID = projectID
	}
	if r.Dest != "" {
		if r.destTable, err = base.NewTableReference(r.Dest); err != nil {
			return err
		}
	}
	if r.destTable != nil {
		if r.destTable.ProjectId == "" {
			r.destTable.ProjectId = projectID
		}
	}
	return nil
}

func (r *QueryRequest) Validate() error {
	if r.SQL == "" {
		return fmt.Errorf("SQL was empty")
	}
	return nil
}

//NewQueryRequest creates a new query request
func NewQueryRequest(SQL string, dest *bigquery.TableReference, finally *task.Actions) *QueryRequest {
	result := &QueryRequest{
		SQL:       SQL,
		destTable: dest,
	}
	if dest != nil {
		result.Dest = dest.DatasetId + "." + dest.TableId
	}
	if finally != nil {
		result.Actions = *finally
	}
	return result
}