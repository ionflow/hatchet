// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: workers.sql

package dbsqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createWorker = `-- name: CreateWorker :one
INSERT INTO "Worker" (
    "id",
    "createdAt",
    "updatedAt",
    "tenantId",
    "name",
    "dispatcherId",
    "maxRuns"
) VALUES (
    gen_random_uuid(),
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    $1::uuid,
    $2::text,
    $3::uuid,
    $4::int
) RETURNING id, "createdAt", "updatedAt", "deletedAt", "tenantId", "lastHeartbeatAt", name, "dispatcherId", "maxRuns", "isActive", "lastListenerEstablished", "isPaused"
`

type CreateWorkerParams struct {
	Tenantid     pgtype.UUID `json:"tenantid"`
	Name         string      `json:"name"`
	Dispatcherid pgtype.UUID `json:"dispatcherid"`
	MaxRuns      pgtype.Int4 `json:"maxRuns"`
}

func (q *Queries) CreateWorker(ctx context.Context, db DBTX, arg CreateWorkerParams) (*Worker, error) {
	row := db.QueryRow(ctx, createWorker,
		arg.Tenantid,
		arg.Name,
		arg.Dispatcherid,
		arg.MaxRuns,
	)
	var i Worker
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.TenantId,
		&i.LastHeartbeatAt,
		&i.Name,
		&i.DispatcherId,
		&i.MaxRuns,
		&i.IsActive,
		&i.LastListenerEstablished,
		&i.IsPaused,
	)
	return &i, err
}

const deleteWorker = `-- name: DeleteWorker :one
DELETE FROM
  "Worker"
WHERE
  "id" = $1::uuid
RETURNING id, "createdAt", "updatedAt", "deletedAt", "tenantId", "lastHeartbeatAt", name, "dispatcherId", "maxRuns", "isActive", "lastListenerEstablished", "isPaused"
`

func (q *Queries) DeleteWorker(ctx context.Context, db DBTX, id pgtype.UUID) (*Worker, error) {
	row := db.QueryRow(ctx, deleteWorker, id)
	var i Worker
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.TenantId,
		&i.LastHeartbeatAt,
		&i.Name,
		&i.DispatcherId,
		&i.MaxRuns,
		&i.IsActive,
		&i.LastListenerEstablished,
		&i.IsPaused,
	)
	return &i, err
}

const getWorkerForEngine = `-- name: GetWorkerForEngine :one
SELECT
    w."id" AS "id",
    w."tenantId" AS "tenantId",
    w."dispatcherId" AS "dispatcherId",
    d."lastHeartbeatAt" AS "dispatcherLastHeartbeatAt",
    w."isActive" AS "isActive",
    w."lastListenerEstablished" AS "lastListenerEstablished"
FROM
    "Worker" w
LEFT JOIN
    "Dispatcher" d ON w."dispatcherId" = d."id"
WHERE
    w."tenantId" = $1
    AND w."id" = $2
`

type GetWorkerForEngineParams struct {
	Tenantid pgtype.UUID `json:"tenantid"`
	ID       pgtype.UUID `json:"id"`
}

type GetWorkerForEngineRow struct {
	ID                        pgtype.UUID      `json:"id"`
	TenantId                  pgtype.UUID      `json:"tenantId"`
	DispatcherId              pgtype.UUID      `json:"dispatcherId"`
	DispatcherLastHeartbeatAt pgtype.Timestamp `json:"dispatcherLastHeartbeatAt"`
	IsActive                  bool             `json:"isActive"`
	LastListenerEstablished   pgtype.Timestamp `json:"lastListenerEstablished"`
}

func (q *Queries) GetWorkerForEngine(ctx context.Context, db DBTX, arg GetWorkerForEngineParams) (*GetWorkerForEngineRow, error) {
	row := db.QueryRow(ctx, getWorkerForEngine, arg.Tenantid, arg.ID)
	var i GetWorkerForEngineRow
	err := row.Scan(
		&i.ID,
		&i.TenantId,
		&i.DispatcherId,
		&i.DispatcherLastHeartbeatAt,
		&i.IsActive,
		&i.LastListenerEstablished,
	)
	return &i, err
}

const linkActionsToWorker = `-- name: LinkActionsToWorker :exec
INSERT INTO "_ActionToWorker" (
    "A",
    "B"
) SELECT
    unnest($1::uuid[]),
    $2::uuid
ON CONFLICT DO NOTHING
`

type LinkActionsToWorkerParams struct {
	Actionids []pgtype.UUID `json:"actionids"`
	Workerid  pgtype.UUID   `json:"workerid"`
}

func (q *Queries) LinkActionsToWorker(ctx context.Context, db DBTX, arg LinkActionsToWorkerParams) error {
	_, err := db.Exec(ctx, linkActionsToWorker, arg.Actionids, arg.Workerid)
	return err
}

const linkServicesToWorker = `-- name: LinkServicesToWorker :exec
INSERT INTO "_ServiceToWorker" (
    "A",
    "B"
)
VALUES (
    unnest($1::uuid[]),
    $2::uuid
)
ON CONFLICT DO NOTHING
`

type LinkServicesToWorkerParams struct {
	Services []pgtype.UUID `json:"services"`
	Workerid pgtype.UUID   `json:"workerid"`
}

func (q *Queries) LinkServicesToWorker(ctx context.Context, db DBTX, arg LinkServicesToWorkerParams) error {
	_, err := db.Exec(ctx, linkServicesToWorker, arg.Services, arg.Workerid)
	return err
}

const listWorkerLabels = `-- name: ListWorkerLabels :many
SELECT
    "id",
    "key",
    "intValue",
    "strValue",
    "createdAt",
    "updatedAt"
FROM "WorkerLabel" wl
WHERE wl."workerId" = $1::uuid
`

type ListWorkerLabelsRow struct {
	ID        int64            `json:"id"`
	Key       string           `json:"key"`
	IntValue  pgtype.Int4      `json:"intValue"`
	StrValue  pgtype.Text      `json:"strValue"`
	CreatedAt pgtype.Timestamp `json:"createdAt"`
	UpdatedAt pgtype.Timestamp `json:"updatedAt"`
}

func (q *Queries) ListWorkerLabels(ctx context.Context, db DBTX, workerid pgtype.UUID) ([]*ListWorkerLabelsRow, error) {
	rows, err := db.Query(ctx, listWorkerLabels, workerid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListWorkerLabelsRow
	for rows.Next() {
		var i ListWorkerLabelsRow
		if err := rows.Scan(
			&i.ID,
			&i.Key,
			&i.IntValue,
			&i.StrValue,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listWorkersWithStepCount = `-- name: ListWorkersWithStepCount :many
SELECT
    workers.id, workers."createdAt", workers."updatedAt", workers."deletedAt", workers."tenantId", workers."lastHeartbeatAt", workers.name, workers."dispatcherId", workers."maxRuns", workers."isActive", workers."lastListenerEstablished", workers."isPaused",
    (SELECT COUNT(*) FROM "WorkerSemaphoreSlot" wss WHERE wss."workerId" = workers."id" AND wss."stepRunId" IS NOT NULL) AS "slots"
FROM
    "Worker" workers
WHERE
    workers."tenantId" = $1
    AND (
        $2::text IS NULL OR
        workers."id" IN (
            SELECT "_ActionToWorker"."B"
            FROM "_ActionToWorker"
            INNER JOIN "Action" ON "Action"."id" = "_ActionToWorker"."A"
            WHERE "Action"."tenantId" = $1 AND "Action"."actionId" = $2::text
        )
    )
    AND (
        $3::timestamp IS NULL OR
        workers."lastHeartbeatAt" > $3::timestamp
    )
    AND (
        $4::boolean IS NULL OR
        workers."maxRuns" IS NULL OR
        ($4::boolean AND workers."maxRuns" > (
            SELECT COUNT(*)
            FROM "StepRun" srs
            WHERE srs."workerId" = workers."id" AND srs."status" = 'RUNNING'
        ))
    )
GROUP BY
    workers."id"
`

type ListWorkersWithStepCountParams struct {
	Tenantid           pgtype.UUID      `json:"tenantid"`
	ActionId           pgtype.Text      `json:"actionId"`
	LastHeartbeatAfter pgtype.Timestamp `json:"lastHeartbeatAfter"`
	Assignable         pgtype.Bool      `json:"assignable"`
}

type ListWorkersWithStepCountRow struct {
	Worker Worker `json:"worker"`
	Slots  int64  `json:"slots"`
}

func (q *Queries) ListWorkersWithStepCount(ctx context.Context, db DBTX, arg ListWorkersWithStepCountParams) ([]*ListWorkersWithStepCountRow, error) {
	rows, err := db.Query(ctx, listWorkersWithStepCount,
		arg.Tenantid,
		arg.ActionId,
		arg.LastHeartbeatAfter,
		arg.Assignable,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListWorkersWithStepCountRow
	for rows.Next() {
		var i ListWorkersWithStepCountRow
		if err := rows.Scan(
			&i.Worker.ID,
			&i.Worker.CreatedAt,
			&i.Worker.UpdatedAt,
			&i.Worker.DeletedAt,
			&i.Worker.TenantId,
			&i.Worker.LastHeartbeatAt,
			&i.Worker.Name,
			&i.Worker.DispatcherId,
			&i.Worker.MaxRuns,
			&i.Worker.IsActive,
			&i.Worker.LastListenerEstablished,
			&i.Worker.IsPaused,
			&i.Slots,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const resolveWorkerSemaphoreSlots = `-- name: ResolveWorkerSemaphoreSlots :one
WITH to_count AS (
    SELECT wss."id"
    FROM "WorkerSemaphoreSlot" wss
    JOIN "StepRun" sr ON wss."stepRunId" = sr."id"
        AND sr."status" NOT IN ('RUNNING', 'ASSIGNED')
        AND sr."tenantId" = $1::uuid
    ORDER BY RANDOM()
    LIMIT 11
    FOR UPDATE SKIP LOCKED
),
to_resolve AS (
    SELECT id FROM to_count LIMIT 10
),
update_result AS (
    UPDATE "WorkerSemaphoreSlot" wss
    SET "stepRunId" = null
    WHERE wss."id" IN (SELECT "id" FROM to_resolve)
    RETURNING wss."id"
)
SELECT
	CASE
		WHEN COUNT(*) > 0 THEN TRUE
		ELSE FALSE
	END AS "hasResolved",
	CASE
		WHEN COUNT(*) > 10 THEN TRUE
		ELSE FALSE
	END AS "hasMore"
FROM to_count
`

type ResolveWorkerSemaphoreSlotsRow struct {
	HasResolved bool `json:"hasResolved"`
	HasMore     bool `json:"hasMore"`
}

func (q *Queries) ResolveWorkerSemaphoreSlots(ctx context.Context, db DBTX, tenantid pgtype.UUID) (*ResolveWorkerSemaphoreSlotsRow, error) {
	row := db.QueryRow(ctx, resolveWorkerSemaphoreSlots, tenantid)
	var i ResolveWorkerSemaphoreSlotsRow
	err := row.Scan(&i.HasResolved, &i.HasMore)
	return &i, err
}

const stubWorkerSemaphoreSlots = `-- name: StubWorkerSemaphoreSlots :exec
INSERT INTO "WorkerSemaphoreSlot" ("id", "workerId")
SELECT gen_random_uuid(), $1::uuid
FROM generate_series(1, $2::int)
`

type StubWorkerSemaphoreSlotsParams struct {
	Workerid pgtype.UUID `json:"workerid"`
	MaxRuns  pgtype.Int4 `json:"maxRuns"`
}

func (q *Queries) StubWorkerSemaphoreSlots(ctx context.Context, db DBTX, arg StubWorkerSemaphoreSlotsParams) error {
	_, err := db.Exec(ctx, stubWorkerSemaphoreSlots, arg.Workerid, arg.MaxRuns)
	return err
}

const updateWorker = `-- name: UpdateWorker :one
UPDATE
    "Worker"
SET
    "updatedAt" = CURRENT_TIMESTAMP,
    "dispatcherId" = coalesce($1::uuid, "dispatcherId"),
    "maxRuns" = coalesce($2::int, "maxRuns"),
    "lastHeartbeatAt" = coalesce($3::timestamp, "lastHeartbeatAt"),
    "isActive" = coalesce($4::boolean, "isActive"),
    "isPaused" = coalesce($5::boolean, "isPaused")
WHERE
    "id" = $6::uuid
RETURNING id, "createdAt", "updatedAt", "deletedAt", "tenantId", "lastHeartbeatAt", name, "dispatcherId", "maxRuns", "isActive", "lastListenerEstablished", "isPaused"
`

type UpdateWorkerParams struct {
	DispatcherId    pgtype.UUID      `json:"dispatcherId"`
	MaxRuns         pgtype.Int4      `json:"maxRuns"`
	LastHeartbeatAt pgtype.Timestamp `json:"lastHeartbeatAt"`
	IsActive        pgtype.Bool      `json:"isActive"`
	IsPaused        pgtype.Bool      `json:"isPaused"`
	ID              pgtype.UUID      `json:"id"`
}

func (q *Queries) UpdateWorker(ctx context.Context, db DBTX, arg UpdateWorkerParams) (*Worker, error) {
	row := db.QueryRow(ctx, updateWorker,
		arg.DispatcherId,
		arg.MaxRuns,
		arg.LastHeartbeatAt,
		arg.IsActive,
		arg.IsPaused,
		arg.ID,
	)
	var i Worker
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.TenantId,
		&i.LastHeartbeatAt,
		&i.Name,
		&i.DispatcherId,
		&i.MaxRuns,
		&i.IsActive,
		&i.LastListenerEstablished,
		&i.IsPaused,
	)
	return &i, err
}

const updateWorkerActiveStatus = `-- name: UpdateWorkerActiveStatus :one
UPDATE "Worker"
SET
    "isActive" = $1::boolean,
    "lastListenerEstablished" = $2::timestamp
WHERE
    "id" = $3::uuid
    AND (
        "lastListenerEstablished" IS NULL
        OR "lastListenerEstablished" <= $2::timestamp
        )
RETURNING id, "createdAt", "updatedAt", "deletedAt", "tenantId", "lastHeartbeatAt", name, "dispatcherId", "maxRuns", "isActive", "lastListenerEstablished", "isPaused"
`

type UpdateWorkerActiveStatusParams struct {
	Isactive                bool             `json:"isactive"`
	LastListenerEstablished pgtype.Timestamp `json:"lastListenerEstablished"`
	ID                      pgtype.UUID      `json:"id"`
}

func (q *Queries) UpdateWorkerActiveStatus(ctx context.Context, db DBTX, arg UpdateWorkerActiveStatusParams) (*Worker, error) {
	row := db.QueryRow(ctx, updateWorkerActiveStatus, arg.Isactive, arg.LastListenerEstablished, arg.ID)
	var i Worker
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.TenantId,
		&i.LastHeartbeatAt,
		&i.Name,
		&i.DispatcherId,
		&i.MaxRuns,
		&i.IsActive,
		&i.LastListenerEstablished,
		&i.IsPaused,
	)
	return &i, err
}

const updateWorkerHeartbeat = `-- name: UpdateWorkerHeartbeat :one
WITH to_update AS (
    SELECT
        "id"
    FROM
        "Worker"
    WHERE
        "id" = $2::uuid
    FOR UPDATE SKIP LOCKED
)
UPDATE
    "Worker"
SET
    "updatedAt" = CURRENT_TIMESTAMP,
    "lastHeartbeatAt" = $1::timestamp
WHERE
    "id" IN (SELECT "id" FROM to_update)
RETURNING id, "createdAt", "updatedAt", "deletedAt", "tenantId", "lastHeartbeatAt", name, "dispatcherId", "maxRuns", "isActive", "lastListenerEstablished", "isPaused"
`

type UpdateWorkerHeartbeatParams struct {
	LastHeartbeatAt pgtype.Timestamp `json:"lastHeartbeatAt"`
	ID              pgtype.UUID      `json:"id"`
}

func (q *Queries) UpdateWorkerHeartbeat(ctx context.Context, db DBTX, arg UpdateWorkerHeartbeatParams) (*Worker, error) {
	row := db.QueryRow(ctx, updateWorkerHeartbeat, arg.LastHeartbeatAt, arg.ID)
	var i Worker
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.TenantId,
		&i.LastHeartbeatAt,
		&i.Name,
		&i.DispatcherId,
		&i.MaxRuns,
		&i.IsActive,
		&i.LastListenerEstablished,
		&i.IsPaused,
	)
	return &i, err
}

const updateWorkersByName = `-- name: UpdateWorkersByName :many
UPDATE "Worker"
SET "isActive" = $1::boolean
WHERE
  "tenantId" = $2::uuid AND
  "name" = $3::text
RETURNING id, "createdAt", "updatedAt", "deletedAt", "tenantId", "lastHeartbeatAt", name, "dispatcherId", "maxRuns", "isActive", "lastListenerEstablished", "isPaused"
`

type UpdateWorkersByNameParams struct {
	Isactive bool        `json:"isactive"`
	Tenantid pgtype.UUID `json:"tenantid"`
	Name     string      `json:"name"`
}

func (q *Queries) UpdateWorkersByName(ctx context.Context, db DBTX, arg UpdateWorkersByNameParams) ([]*Worker, error) {
	rows, err := db.Query(ctx, updateWorkersByName, arg.Isactive, arg.Tenantid, arg.Name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Worker
	for rows.Next() {
		var i Worker
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.TenantId,
			&i.LastHeartbeatAt,
			&i.Name,
			&i.DispatcherId,
			&i.MaxRuns,
			&i.IsActive,
			&i.LastListenerEstablished,
			&i.IsPaused,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const upsertService = `-- name: UpsertService :one
INSERT INTO "Service" (
    "id",
    "createdAt",
    "updatedAt",
    "name",
    "tenantId"
)
VALUES (
    gen_random_uuid(),
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    $1::text,
    $2::uuid
)
ON CONFLICT ("tenantId", "name") DO UPDATE
SET
    "updatedAt" = CURRENT_TIMESTAMP
WHERE
    "Service"."tenantId" = $2 AND "Service"."name" = $1::text
RETURNING id, "createdAt", "updatedAt", "deletedAt", name, description, "tenantId"
`

type UpsertServiceParams struct {
	Name     string      `json:"name"`
	Tenantid pgtype.UUID `json:"tenantid"`
}

func (q *Queries) UpsertService(ctx context.Context, db DBTX, arg UpsertServiceParams) (*Service, error) {
	row := db.QueryRow(ctx, upsertService, arg.Name, arg.Tenantid)
	var i Service
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Name,
		&i.Description,
		&i.TenantId,
	)
	return &i, err
}

const upsertWorkerLabel = `-- name: UpsertWorkerLabel :one
INSERT INTO "WorkerLabel" (
    "createdAt",
    "updatedAt",
    "workerId",
    "key",
    "intValue",
    "strValue"
) VALUES (
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    $1::uuid,
    $2::text,
    $3::int,
    $4::text
) ON CONFLICT ("workerId", "key") DO UPDATE
SET
    "updatedAt" = CURRENT_TIMESTAMP,
    "intValue" = $3::int,
    "strValue" = $4::text
RETURNING id, "createdAt", "updatedAt", "workerId", key, "strValue", "intValue"
`

type UpsertWorkerLabelParams struct {
	Workerid pgtype.UUID `json:"workerid"`
	Key      string      `json:"key"`
	IntValue pgtype.Int4 `json:"intValue"`
	StrValue pgtype.Text `json:"strValue"`
}

func (q *Queries) UpsertWorkerLabel(ctx context.Context, db DBTX, arg UpsertWorkerLabelParams) (*WorkerLabel, error) {
	row := db.QueryRow(ctx, upsertWorkerLabel,
		arg.Workerid,
		arg.Key,
		arg.IntValue,
		arg.StrValue,
	)
	var i WorkerLabel
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.WorkerId,
		&i.Key,
		&i.StrValue,
		&i.IntValue,
	)
	return &i, err
}
