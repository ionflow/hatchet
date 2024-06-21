package prisma

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"github.com/hatchet-dev/hatchet/pkg/repository"
	"github.com/hatchet-dev/hatchet/pkg/repository/prisma/dbsqlc"
	"github.com/hatchet-dev/hatchet/pkg/repository/prisma/sqlchelpers"
	"github.com/hatchet-dev/hatchet/pkg/validator"
)

type webhookWorkerEngineRepository struct {
	pool    *pgxpool.Pool
	v       validator.Validator
	queries *dbsqlc.Queries
	l       *zerolog.Logger
}

func NewWebhookWorkerEngineRepository(pool *pgxpool.Pool, v validator.Validator, l *zerolog.Logger) repository.WebhookWorkerEngineRepository {
	queries := dbsqlc.New()

	return &webhookWorkerEngineRepository{
		pool:    pool,
		v:       v,
		queries: queries,
		l:       l,
	}
}

func (r *webhookWorkerEngineRepository) ListWebhookWorkers(ctx context.Context, tenantId string) ([]*dbsqlc.WebhookWorker, error) {
	return r.queries.ListWebhookWorkers(ctx, r.pool, sqlchelpers.UUIDFromStr(tenantId))
}

func (r *webhookWorkerEngineRepository) UpsertWebhookWorker(ctx context.Context, opts *repository.UpsertWebhookWorkerOpts) (*dbsqlc.WebhookWorker, error) {
	if err := r.v.Validate(opts); err != nil {
		return nil, err
	}

	params := dbsqlc.UpsertWebhookWorkerParams{
		Tenantid: sqlchelpers.UUIDFromStr(opts.TenantId),
		Name:     sqlchelpers.TextFromStr(opts.Name),
		Secret:   sqlchelpers.TextFromStr(opts.Secret),
		Url:      sqlchelpers.TextFromStr(opts.URL),
	}

	if opts.TokenID != nil {
		params.TokenId = sqlchelpers.UUIDFromStr(*opts.TokenID)
	}

	if opts.TokenValue != nil {
		params.TokenValue = sqlchelpers.TextFromStr(*opts.TokenValue)
	}

	return r.queries.UpsertWebhookWorker(ctx, r.pool, params)
}
