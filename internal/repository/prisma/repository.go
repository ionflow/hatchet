package prisma

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/repository"
	"github.com/hatchet-dev/hatchet/internal/repository/cache"
	"github.com/hatchet-dev/hatchet/internal/repository/metered"
	"github.com/hatchet-dev/hatchet/internal/repository/prisma/db"
	"github.com/hatchet-dev/hatchet/internal/validator"
)

type apiRepository struct {
	apiToken       repository.APITokenRepository
	event          repository.EventAPIRepository
	log            repository.LogsAPIRepository
	tenant         repository.TenantAPIRepository
	tenantAlerting repository.TenantAlertingAPIRepository
	tenantInvite   repository.TenantInviteRepository
	workflow       repository.WorkflowAPIRepository
	workflowRun    repository.WorkflowRunAPIRepository
	jobRun         repository.JobRunAPIRepository
	stepRun        repository.StepRunAPIRepository
	github         repository.GithubRepository
	step           repository.StepRepository
	slack          repository.SlackRepository
	sns            repository.SNSRepository
	worker         repository.WorkerAPIRepository
	userSession    repository.UserSessionRepository
	user           repository.UserRepository
	health         repository.HealthRepository
	webhookWorker  repository.WebhookWorkerRepository
}

type PrismaRepositoryOpt func(*PrismaRepositoryOpts)

type PrismaRepositoryOpts struct {
	v       validator.Validator
	l       *zerolog.Logger
	cache   cache.Cacheable
	metered *metered.Metered
}

func defaultPrismaRepositoryOpts() *PrismaRepositoryOpts {
	return &PrismaRepositoryOpts{
		v: validator.NewDefaultValidator(),
	}
}

func WithValidator(v validator.Validator) PrismaRepositoryOpt {
	return func(opts *PrismaRepositoryOpts) {
		opts.v = v
	}
}

func WithLogger(l *zerolog.Logger) PrismaRepositoryOpt {
	return func(opts *PrismaRepositoryOpts) {
		opts.l = l
	}
}

func WithCache(cache cache.Cacheable) PrismaRepositoryOpt {
	return func(opts *PrismaRepositoryOpts) {
		opts.cache = cache
	}
}

func WithMetered(metered *metered.Metered) PrismaRepositoryOpt {
	return func(opts *PrismaRepositoryOpts) {
		opts.metered = metered
	}
}

func NewAPIRepository(client *db.PrismaClient, pool *pgxpool.Pool, fs ...PrismaRepositoryOpt) repository.APIRepository {
	opts := defaultPrismaRepositoryOpts()

	for _, f := range fs {
		f(opts)
	}

	newLogger := opts.l.With().Str("service", "database").Logger()
	opts.l = &newLogger

	if opts.cache == nil {
		opts.cache = cache.New(1 * time.Millisecond)
	}

	return &apiRepository{
		apiToken:       NewAPITokenRepository(client, opts.v, opts.cache),
		event:          NewEventAPIRepository(client, pool, opts.v, opts.l),
		log:            NewLogAPIRepository(pool, opts.v, opts.l),
		tenant:         NewTenantAPIRepository(client, opts.v, opts.cache),
		tenantAlerting: NewTenantAlertingAPIRepository(client, opts.v, opts.cache),
		tenantInvite:   NewTenantInviteRepository(client, opts.v),
		workflow:       NewWorkflowRepository(client, pool, opts.v, opts.l),
		workflowRun:    NewWorkflowRunRepository(client, pool, opts.v, opts.l, opts.metered),
		jobRun:         NewJobRunAPIRepository(client, pool, opts.v, opts.l),
		stepRun:        NewStepRunAPIRepository(client, pool, opts.v, opts.l),
		github:         NewGithubRepository(client, opts.v),
		step:           NewStepRepository(client, opts.v),
		slack:          NewSlackRepository(client, opts.v),
		sns:            NewSNSRepository(client, opts.v),
		worker:         NewWorkerAPIRepository(client, pool, opts.v, opts.l, opts.metered),
		userSession:    NewUserSessionRepository(client, opts.v),
		user:           NewUserRepository(client, opts.v),
		health:         NewHealthAPIRepository(client, pool),
		webhookWorker:  NewWebhookWorkerRepository(client, opts.v, opts.l),
	}
}

func (r *apiRepository) Health() repository.HealthRepository {
	return r.health
}

func (r *apiRepository) APIToken() repository.APITokenRepository {
	return r.apiToken
}

func (r *apiRepository) Event() repository.EventAPIRepository {
	return r.event
}

func (r *apiRepository) Log() repository.LogsAPIRepository {
	return r.log
}

func (r *apiRepository) Tenant() repository.TenantAPIRepository {
	return r.tenant
}

func (r *apiRepository) TenantAlertingSettings() repository.TenantAlertingAPIRepository {
	return r.tenantAlerting
}

func (r *apiRepository) TenantInvite() repository.TenantInviteRepository {
	return r.tenantInvite
}

func (r *apiRepository) Workflow() repository.WorkflowAPIRepository {
	return r.workflow
}

func (r *apiRepository) WorkflowRun() repository.WorkflowRunAPIRepository {
	return r.workflowRun
}

func (r *apiRepository) JobRun() repository.JobRunAPIRepository {
	return r.jobRun
}

func (r *apiRepository) StepRun() repository.StepRunAPIRepository {
	return r.stepRun
}

func (r *apiRepository) Slack() repository.SlackRepository {
	return r.slack
}

func (r *apiRepository) SNS() repository.SNSRepository {
	return r.sns
}

func (r *apiRepository) Github() repository.GithubRepository {
	return r.github
}

func (r *apiRepository) Step() repository.StepRepository {
	return r.step
}

func (r *apiRepository) Worker() repository.WorkerAPIRepository {
	return r.worker
}

func (r *apiRepository) UserSession() repository.UserSessionRepository {
	return r.userSession
}

func (r *apiRepository) User() repository.UserRepository {
	return r.user
}

func (r *apiRepository) WebhookWorker() repository.WebhookWorkerRepository {
	return r.webhookWorker
}

type engineRepository struct {
	health         repository.HealthRepository
	apiToken       repository.EngineTokenRepository
	dispatcher     repository.DispatcherEngineRepository
	event          repository.EventEngineRepository
	getGroupKeyRun repository.GetGroupKeyRunEngineRepository
	jobRun         repository.JobRunEngineRepository
	stepRun        repository.StepRunEngineRepository
	tenant         repository.TenantEngineRepository
	tenantAlerting repository.TenantAlertingEngineRepository
	ticker         repository.TickerEngineRepository
	worker         repository.WorkerEngineRepository
	workflow       repository.WorkflowEngineRepository
	workflowRun    repository.WorkflowRunEngineRepository
	streamEvent    repository.StreamEventsEngineRepository
	log            repository.LogsEngineRepository
	rateLimit      repository.RateLimitEngineRepository
}

func (r *engineRepository) Health() repository.HealthRepository {
	return r.health
}

func (r *engineRepository) APIToken() repository.EngineTokenRepository {
	return r.apiToken
}

func (r *engineRepository) Dispatcher() repository.DispatcherEngineRepository {
	return r.dispatcher
}

func (r *engineRepository) Event() repository.EventEngineRepository {
	return r.event
}

func (r *engineRepository) GetGroupKeyRun() repository.GetGroupKeyRunEngineRepository {
	return r.getGroupKeyRun
}

func (r *engineRepository) JobRun() repository.JobRunEngineRepository {
	return r.jobRun
}

func (r *engineRepository) StepRun() repository.StepRunEngineRepository {
	return r.stepRun
}

func (r *engineRepository) Tenant() repository.TenantEngineRepository {
	return r.tenant
}

func (r *engineRepository) TenantAlertingSettings() repository.TenantAlertingEngineRepository {
	return r.tenantAlerting
}

func (r *engineRepository) Ticker() repository.TickerEngineRepository {
	return r.ticker
}

func (r *engineRepository) Worker() repository.WorkerEngineRepository {
	return r.worker
}

func (r *engineRepository) Workflow() repository.WorkflowEngineRepository {
	return r.workflow
}

func (r *engineRepository) WorkflowRun() repository.WorkflowRunEngineRepository {
	return r.workflowRun
}

func (r *engineRepository) StreamEvent() repository.StreamEventsEngineRepository {
	return r.streamEvent
}

func (r *engineRepository) Log() repository.LogsEngineRepository {
	return r.log
}

func (r *engineRepository) RateLimit() repository.RateLimitEngineRepository {
	return r.rateLimit
}

func NewEngineRepository(pool *pgxpool.Pool, fs ...PrismaRepositoryOpt) repository.EngineRepository {
	opts := defaultPrismaRepositoryOpts()

	for _, f := range fs {
		f(opts)
	}

	newLogger := opts.l.With().Str("service", "database").Logger()
	opts.l = &newLogger

	if opts.cache == nil {
		opts.cache = cache.New(1 * time.Millisecond)
	}

	return &engineRepository{
		health:         NewHealthEngineRepository(pool),
		apiToken:       NewEngineTokenRepository(pool, opts.v, opts.l, opts.cache),
		dispatcher:     NewDispatcherRepository(pool, opts.v, opts.l),
		event:          NewEventEngineRepository(pool, opts.v, opts.l, opts.metered),
		getGroupKeyRun: NewGetGroupKeyRunRepository(pool, opts.v, opts.l),
		jobRun:         NewJobRunEngineRepository(pool, opts.v, opts.l),
		stepRun:        NewStepRunEngineRepository(pool, opts.v, opts.l),
		tenant:         NewTenantEngineRepository(pool, opts.v, opts.l, opts.cache),
		tenantAlerting: NewTenantAlertingEngineRepository(pool, opts.v, opts.l, opts.cache),
		ticker:         NewTickerRepository(pool, opts.v, opts.l),
		worker:         NewWorkerEngineRepository(pool, opts.v, opts.l, opts.metered),
		workflow:       NewWorkflowEngineRepository(pool, opts.v, opts.l, opts.metered),
		workflowRun:    NewWorkflowRunEngineRepository(pool, opts.v, opts.l, opts.metered),
		streamEvent:    NewStreamEventsEngineRepository(pool, opts.v, opts.l),
		log:            NewLogEngineRepository(pool, opts.v, opts.l),
		rateLimit:      NewRateLimitEngineRepository(pool, opts.v, opts.l),
	}
}

type entitlementRepository struct {
	tenantLimit repository.TenantLimitRepository
}

func (r *entitlementRepository) TenantLimit() repository.TenantLimitRepository {
	return r.tenantLimit
}
func NewEntitlementRepository(pool *pgxpool.Pool, s *server.ConfigFileRuntime, fs ...PrismaRepositoryOpt) repository.EntitlementsRepository {
	opts := defaultPrismaRepositoryOpts()

	for _, f := range fs {
		f(opts)
	}

	newLogger := opts.l.With().Str("service", "database").Logger()
	opts.l = &newLogger

	if opts.cache == nil {
		opts.cache = cache.New(1 * time.Millisecond)
	}

	return &entitlementRepository{
		tenantLimit: NewTenantLimitRepository(pool, opts.v, opts.l, s),
	}
}
