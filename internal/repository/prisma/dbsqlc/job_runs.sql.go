// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: job_runs.sql

package dbsqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getJobRunByWorkflowRunIdAndJobId = `-- name: GetJobRunByWorkflowRunIdAndJobId :one
SELECT
    "id",
    "jobId",
    "status"
FROM
    "JobRun" jr
WHERE
    jr."tenantId" = $1::uuid
    AND jr."workflowRunId" = $2::uuid
    AND jr."jobId" = $3::uuid
`

type GetJobRunByWorkflowRunIdAndJobIdParams struct {
	Tenantid      pgtype.UUID `json:"tenantid"`
	Workflowrunid pgtype.UUID `json:"workflowrunid"`
	Jobid         pgtype.UUID `json:"jobid"`
}

type GetJobRunByWorkflowRunIdAndJobIdRow struct {
	ID     pgtype.UUID  `json:"id"`
	JobId  pgtype.UUID  `json:"jobId"`
	Status JobRunStatus `json:"status"`
}

func (q *Queries) GetJobRunByWorkflowRunIdAndJobId(ctx context.Context, db DBTX, arg GetJobRunByWorkflowRunIdAndJobIdParams) (*GetJobRunByWorkflowRunIdAndJobIdRow, error) {
	row := db.QueryRow(ctx, getJobRunByWorkflowRunIdAndJobId, arg.Tenantid, arg.Workflowrunid, arg.Jobid)
	var i GetJobRunByWorkflowRunIdAndJobIdRow
	err := row.Scan(&i.ID, &i.JobId, &i.Status)
	return &i, err
}

const listJobRunsForWorkflowRun = `-- name: ListJobRunsForWorkflowRun :many
SELECT
    "id",
    "jobId"
FROM
    "JobRun" jr
WHERE
    jr."workflowRunId" = $1::uuid
`

type ListJobRunsForWorkflowRunRow struct {
	ID    pgtype.UUID `json:"id"`
	JobId pgtype.UUID `json:"jobId"`
}

func (q *Queries) ListJobRunsForWorkflowRun(ctx context.Context, db DBTX, workflowrunid pgtype.UUID) ([]*ListJobRunsForWorkflowRunRow, error) {
	rows, err := db.Query(ctx, listJobRunsForWorkflowRun, workflowrunid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListJobRunsForWorkflowRunRow
	for rows.Next() {
		var i ListJobRunsForWorkflowRunRow
		if err := rows.Scan(&i.ID, &i.JobId); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const resolveJobRunStatus = `-- name: ResolveJobRunStatus :one
WITH stepRuns AS (
    SELECT sum(case when runs."status" IN ('PENDING', 'PENDING_ASSIGNMENT') then 1 else 0 end) AS pendingRuns,
        sum(case when runs."status" IN ('RUNNING', 'ASSIGNED') then 1 else 0 end) AS runningRuns,
        sum(case when runs."status" = 'SUCCEEDED' then 1 else 0 end) AS succeededRuns,
        sum(case when runs."status" = 'FAILED' then 1 else 0 end) AS failedRuns,
        sum(case when runs."status" = 'CANCELLED' then 1 else 0 end) AS cancelledRuns
    FROM "StepRun" as runs
    WHERE
        "jobRunId" = (
            SELECT "jobRunId"
            FROM "StepRun"
            WHERE "id" = $1::uuid
        ) AND
        "tenantId" = $2::uuid
)
UPDATE "JobRun"
SET "status" = CASE
    -- Final states are final, cannot be updated
    WHEN "status" IN ('SUCCEEDED', 'FAILED', 'CANCELLED') THEN "status"
    -- NOTE: Order of the following conditions is important
    -- When one step run is running, then the job is running
    WHEN (s.runningRuns > 0 OR s.pendingRuns > 0) THEN 'RUNNING'
    -- When one step run has failed, then the job is failed
    WHEN s.failedRuns > 0 THEN 'FAILED'
    -- When one step run has been cancelled, then the job is cancelled
    WHEN s.cancelledRuns > 0 THEN 'CANCELLED'
    -- When no step runs exist that are not succeeded, then the job is succeeded
    WHEN s.succeededRuns > 0 AND s.pendingRuns = 0 AND s.runningRuns = 0 AND s.failedRuns = 0 AND s.cancelledRuns = 0 THEN 'SUCCEEDED'
    ELSE "status"
END, "finishedAt" = CASE
    -- Final states are final, cannot be updated
    WHEN "finishedAt" IS NOT NULL THEN "finishedAt"
    WHEN s.runningRuns > 0 THEN NULL
    -- When one step run has failed or been cancelled, then the job is finished
    WHEN s.failedRuns > 0 OR s.cancelledRuns > 0 THEN NOW()
    -- When no step runs exist that are not succeeded, then the job is finished
    WHEN s.succeededRuns > 0 AND s.pendingRuns = 0 AND s.runningRuns = 0 AND s.failedRuns = 0 AND s.cancelledRuns = 0 THEN NOW()
    ELSE "finishedAt"
END, "startedAt" = CASE
    -- Started at is final, cannot be changed
    WHEN "startedAt" IS NOT NULL THEN "startedAt"
    -- If steps are running (or have finished), then set the started at time
    WHEN s.runningRuns > 0 OR s.succeededRuns > 0 OR s.failedRuns > 0 AND s.cancelledRuns > 0 THEN NOW()
    ELSE "startedAt"
END
FROM stepRuns s
WHERE "id" = (
    SELECT "jobRunId"
    FROM "StepRun"
    WHERE "id" = $1::uuid
) AND "tenantId" = $2::uuid
RETURNING "JobRun".id, "JobRun"."createdAt", "JobRun"."updatedAt", "JobRun"."deletedAt", "JobRun"."tenantId", "JobRun"."jobId", "JobRun"."tickerId", "JobRun".status, "JobRun".result, "JobRun"."startedAt", "JobRun"."finishedAt", "JobRun"."timeoutAt", "JobRun"."cancelledAt", "JobRun"."cancelledReason", "JobRun"."cancelledError", "JobRun"."workflowRunId"
`

type ResolveJobRunStatusParams struct {
	Steprunid pgtype.UUID `json:"steprunid"`
	Tenantid  pgtype.UUID `json:"tenantid"`
}

func (q *Queries) ResolveJobRunStatus(ctx context.Context, db DBTX, arg ResolveJobRunStatusParams) (*JobRun, error) {
	row := db.QueryRow(ctx, resolveJobRunStatus, arg.Steprunid, arg.Tenantid)
	var i JobRun
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.TenantId,
		&i.JobId,
		&i.TickerId,
		&i.Status,
		&i.Result,
		&i.StartedAt,
		&i.FinishedAt,
		&i.TimeoutAt,
		&i.CancelledAt,
		&i.CancelledReason,
		&i.CancelledError,
		&i.WorkflowRunId,
	)
	return &i, err
}

const updateJobRunLookupDataWithStepRun = `-- name: UpdateJobRunLookupDataWithStepRun :exec
WITH readable_id AS (
    SELECT "readableId"
    FROM "Step"
    WHERE "id" = (
        SELECT "stepId"
        FROM "StepRun"
        WHERE "id" = $2::uuid
    )
)
UPDATE "JobRunLookupData"
SET
    "data" = CASE
        WHEN $1::jsonb IS NULL THEN
            jsonb_set(
                "data",
                '{steps}',
                ("data"->'steps') - (SELECT "readableId" FROM readable_id),
                true
            )
        ELSE
            jsonb_set(
                "data",
                ARRAY['steps', (SELECT "readableId" FROM readable_id)],
                $1::jsonb,
                true
            )
    END,
    "updatedAt" = CURRENT_TIMESTAMP
WHERE
    "jobRunId" = (
        SELECT "jobRunId"
        FROM "StepRun"
        WHERE "id" = $2::uuid
    )
    AND "tenantId" = $3::uuid
`

type UpdateJobRunLookupDataWithStepRunParams struct {
	Jsondata  []byte      `json:"jsondata"`
	Steprunid pgtype.UUID `json:"steprunid"`
	Tenantid  pgtype.UUID `json:"tenantid"`
}

func (q *Queries) UpdateJobRunLookupDataWithStepRun(ctx context.Context, db DBTX, arg UpdateJobRunLookupDataWithStepRunParams) error {
	_, err := db.Exec(ctx, updateJobRunLookupDataWithStepRun, arg.Jsondata, arg.Steprunid, arg.Tenantid)
	return err
}

const updateJobRunStatus = `-- name: UpdateJobRunStatus :one
UPDATE "JobRun"
SET "status" = $1::"JobRunStatus"
WHERE "id" = $2::uuid AND "tenantId" = $3::uuid
RETURNING id, "createdAt", "updatedAt", "deletedAt", "tenantId", "jobId", "tickerId", status, result, "startedAt", "finishedAt", "timeoutAt", "cancelledAt", "cancelledReason", "cancelledError", "workflowRunId"
`

type UpdateJobRunStatusParams struct {
	Status   JobRunStatus `json:"status"`
	ID       pgtype.UUID  `json:"id"`
	Tenantid pgtype.UUID  `json:"tenantid"`
}

func (q *Queries) UpdateJobRunStatus(ctx context.Context, db DBTX, arg UpdateJobRunStatusParams) (*JobRun, error) {
	row := db.QueryRow(ctx, updateJobRunStatus, arg.Status, arg.ID, arg.Tenantid)
	var i JobRun
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.TenantId,
		&i.JobId,
		&i.TickerId,
		&i.Status,
		&i.Result,
		&i.StartedAt,
		&i.FinishedAt,
		&i.TimeoutAt,
		&i.CancelledAt,
		&i.CancelledReason,
		&i.CancelledError,
		&i.WorkflowRunId,
	)
	return &i, err
}

const upsertJobRunLookupData = `-- name: UpsertJobRunLookupData :exec
INSERT INTO "JobRunLookupData" (
    "id",
    "createdAt",
    "updatedAt",
    "deletedAt",
    "jobRunId",
    "tenantId",
    "data"
) VALUES (
    gen_random_uuid(), -- Generates a new UUID for id
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    NULL,
    $1::uuid,
    $2::uuid,
    jsonb_set('{}', $3::text[], $4::jsonb, true)
) ON CONFLICT ("jobRunId", "tenantId") DO UPDATE
SET
    "data" = jsonb_set("JobRunLookupData"."data", $3::text[], $4::jsonb, true),
    "updatedAt" = CURRENT_TIMESTAMP
`

type UpsertJobRunLookupDataParams struct {
	Jobrunid  pgtype.UUID `json:"jobrunid"`
	Tenantid  pgtype.UUID `json:"tenantid"`
	Fieldpath []string    `json:"fieldpath"`
	Jsondata  []byte      `json:"jsondata"`
}

func (q *Queries) UpsertJobRunLookupData(ctx context.Context, db DBTX, arg UpsertJobRunLookupDataParams) error {
	_, err := db.Exec(ctx, upsertJobRunLookupData,
		arg.Jobrunid,
		arg.Tenantid,
		arg.Fieldpath,
		arg.Jsondata,
	)
	return err
}
