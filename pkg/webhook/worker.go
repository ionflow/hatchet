package webhook

import (
	"fmt"
	"os"

	"github.com/hatchet-dev/hatchet/pkg/client"
	"github.com/hatchet-dev/hatchet/pkg/worker"
)

type WebhookWorker struct {
	opts   WorkerOpts
	client client.Client
}

type WorkerOpts struct {
	Token     string
	ID        string
	Secret    string
	URL       string
	TenantID  string
	Actions   []string
	Workflows []string
}

func NewWorker(opts WorkerOpts) (*WebhookWorker, error) {
	// TODO temp hack
	_ = os.Setenv("HATCHET_CLIENT_TOKEN", opts.Token)
	// client.WithToken(opts.Token),
	_ = os.Setenv("HATCHET_CLIENT_TLS_STRATEGY", "none")

	cl, err := client.New()
	if err != nil {
		return nil, fmt.Errorf("could not create client: %w", err)
	}

	return &WebhookWorker{
		opts:   opts,
		client: cl,
	}, nil
}

func (w *WebhookWorker) Start() (func() error, error) {
	r, err := worker.NewWorker(
		worker.WithClient(w.client),
		worker.WithInternalData(w.opts.Actions, w.opts.Workflows),
		worker.WithName("Webhook_"+w.opts.ID),
	)
	if err != nil {
		return nil, fmt.Errorf("could not create webhook worker: %w", err)
	}

	cleanup, err := r.StartWebhook(worker.WebhookWorkerOpts{
		URL:    w.opts.URL,
		Secret: w.opts.Secret,
	})
	if err != nil {
		return nil, fmt.Errorf("could not start webhook worker: %w", err)
	}

	return cleanup, nil
}
