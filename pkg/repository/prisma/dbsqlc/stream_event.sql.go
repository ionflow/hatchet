// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: stream_event.sql

package dbsqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const cleanupStreamEvents = `-- name: CleanupStreamEvents :exec
DELETE FROM "StreamEvent"
WHERE
  -- older than than 5 minutes ago
  "createdAt" < NOW() - INTERVAL '5 minutes'
`

func (q *Queries) CleanupStreamEvents(ctx context.Context, db DBTX) error {
	_, err := db.Exec(ctx, cleanupStreamEvents)
	return err
}

const createStreamEvent = `-- name: CreateStreamEvent :one
INSERT INTO "StreamEvent" (
    "createdAt",
    "tenantId",
    "stepRunId",
    "message",
    "metadata"
)
SELECT
    coalesce($1::timestamp, now()),
    $2::uuid,
    $3::uuid,
    $4::bytea,
    coalesce($5::jsonb, '{}'::jsonb)
FROM "StepRun"
WHERE "StepRun"."id" = $3::uuid
AND "StepRun"."tenantId" = $2::uuid
RETURNING id, "createdAt", "tenantId", "stepRunId", message, metadata
`

type CreateStreamEventParams struct {
	CreatedAt pgtype.Timestamp `json:"createdAt"`
	Tenantid  pgtype.UUID      `json:"tenantid"`
	Steprunid pgtype.UUID      `json:"steprunid"`
	Message   []byte           `json:"message"`
	Metadata  []byte           `json:"metadata"`
}

func (q *Queries) CreateStreamEvent(ctx context.Context, db DBTX, arg CreateStreamEventParams) (*StreamEvent, error) {
	row := db.QueryRow(ctx, createStreamEvent,
		arg.CreatedAt,
		arg.Tenantid,
		arg.Steprunid,
		arg.Message,
		arg.Metadata,
	)
	var i StreamEvent
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.TenantId,
		&i.StepRunId,
		&i.Message,
		&i.Metadata,
	)
	return &i, err
}

const getStreamEvent = `-- name: GetStreamEvent :one
SELECT id, "createdAt", "tenantId", "stepRunId", message, metadata FROM "StreamEvent"
WHERE
  "tenantId" = $1::uuid AND
  "id" = $2::bigint
`

type GetStreamEventParams struct {
	Tenantid pgtype.UUID `json:"tenantid"`
	ID       int64       `json:"id"`
}

func (q *Queries) GetStreamEvent(ctx context.Context, db DBTX, arg GetStreamEventParams) (*StreamEvent, error) {
	row := db.QueryRow(ctx, getStreamEvent, arg.Tenantid, arg.ID)
	var i StreamEvent
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.TenantId,
		&i.StepRunId,
		&i.Message,
		&i.Metadata,
	)
	return &i, err
}

const getStreamEventMeta = `-- name: GetStreamEventMeta :one
SELECT
    jr."workflowRunId" AS "workflowRunId",
    sr."retryCount" AS "retryCount",
    s."retries" as "retries"
FROM "StepRun" sr
JOIN "Step" s ON sr."stepId" = s."id"
JOIN "JobRun" jr ON sr."jobRunId" = jr."id"
WHERE sr."id" = $1::uuid
AND sr."tenantId" = $2::uuid
`

type GetStreamEventMetaParams struct {
	Steprunid pgtype.UUID `json:"steprunid"`
	Tenantid  pgtype.UUID `json:"tenantid"`
}

type GetStreamEventMetaRow struct {
	WorkflowRunId pgtype.UUID `json:"workflowRunId"`
	RetryCount    int32       `json:"retryCount"`
	Retries       int32       `json:"retries"`
}

func (q *Queries) GetStreamEventMeta(ctx context.Context, db DBTX, arg GetStreamEventMetaParams) (*GetStreamEventMetaRow, error) {
	row := db.QueryRow(ctx, getStreamEventMeta, arg.Steprunid, arg.Tenantid)
	var i GetStreamEventMetaRow
	err := row.Scan(&i.WorkflowRunId, &i.RetryCount, &i.Retries)
	return &i, err
}
