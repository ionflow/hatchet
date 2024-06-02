package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/hatchet-dev/hatchet/internal/repository/prisma/db"
	"github.com/hatchet-dev/hatchet/pkg/client"
	"github.com/hatchet-dev/hatchet/pkg/client/rest"
	"github.com/hatchet-dev/hatchet/pkg/worker"
)

func run(c client.Client, job worker.WorkflowJob) error {
	w, err := worker.NewWorker(
		worker.WithClient(
			c,
		),
	)
	if err != nil {
		return fmt.Errorf("error creating worker: %w", err)
	}

	client := db.NewClient()
	if err := client.Connect(); err != nil {
		panic(fmt.Errorf("error connecting to database: %w", err))
	}
	defer client.Disconnect()

	err = w.On(worker.Events("user:create:webhook"), &job)
	if err != nil {
		return fmt.Errorf("error registering webhook workflow: %w", err)
	}

	go func() {
		// create webserver to handle webhook requests
		http.HandleFunc("/webhook", w.WebhookHandler("secret"))

		port := "8741"
		log.Printf("starting webhook server on port %s", port)
		if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	log.Printf("pushing event")

	testEvent := userCreateEvent{
		Username: "echo-test",
		UserID:   "1234",
		Data: map[string]string{
			"test": "test",
		},
	}

	// push an event
	err = c.Event().Push(
		context.Background(),
		"user:create:webhook",
		testEvent,
	)
	if err != nil {
		panic(fmt.Errorf("error pushing event: %w", err))
	}

	// TODO test for assigned status before it is started
	//time.Sleep(2 * time.Second)
	//verifyStepRuns(client, c.TenantId(), db.JobRunStatusRunning, db.StepRunStatusAssigned, nil)

	time.Sleep(5 * time.Second)

	verifyStepRuns(client, c.TenantId(), db.JobRunStatusSucceeded, db.StepRunStatusSucceeded, func(output string) {
		if string(output) != `{"message":"hi from step-one"}` && string(output) != `{"message":"hi from step-two"}` {
			panic(fmt.Errorf("expected step run output to be valid, got %s", output))
		}
	})

	return nil
}

func setup(c client.Client) error {
	tenantId := openapi_types.UUID{}
	if err := tenantId.Scan(c.TenantId()); err != nil {
		return fmt.Errorf("error getting tenant id: %w", err)
	}

	res, err := c.API().WebhookCreate(context.Background(), tenantId, rest.WebhookCreateJSONRequestBody{
		Url: "http://localhost:8741/webhook",
	})
	if err != nil {
		return fmt.Errorf("error creating webhook worker: %w", err)
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("error creating webhook, failed with status code %d", res.StatusCode)
	}

	return nil
}

func verifyStepRuns(client *db.PrismaClient, tenantId string, jobRunStatus db.JobRunStatus, stepRunStatus db.StepRunStatus, check func(string)) {
	events, err := client.Event.FindMany(
		db.Event.TenantID.Equals(tenantId),
		db.Event.Key.Equals("user:create:webhook"),
	).With(
		db.Event.WorkflowRuns.Fetch().With(
			db.WorkflowRunTriggeredBy.Parent.Fetch().With(
				db.WorkflowRun.JobRuns.Fetch().With(
					db.JobRun.StepRuns.Fetch(),
				),
			),
		),
	).Exec(context.Background())
	if err != nil {
		panic(fmt.Errorf("error finding events: %w", err))
	}

	if len(events) == 0 {
		panic(fmt.Errorf("no events found"))
	}

	for _, event := range events {
		if len(event.WorkflowRuns()) == 0 {
			panic(fmt.Errorf("no workflow runs found"))
		}
		for _, workflowRun := range event.WorkflowRuns() {
			if len(workflowRun.Parent().JobRuns()) == 0 {
				panic(fmt.Errorf("no job runs found"))
			}
			for _, jobRuns := range workflowRun.Parent().JobRuns() {
				if jobRuns.Status != jobRunStatus {
					panic(fmt.Errorf("expected job run to be %s, got %s", jobRunStatus, jobRuns.Status))
				}
				for _, stepRun := range jobRuns.StepRuns() {
					if stepRun.Status != stepRunStatus {
						panic(fmt.Errorf("expected step run to be %s, got %s", stepRunStatus, stepRun.Status))
					}
					output, ok := stepRun.Output()
					if check != nil {
						if !ok {
							panic(fmt.Errorf("expected step run to have output, got %+v", stepRun))
						}
						check(string(output))
					}
				}
			}
		}
	}
}
