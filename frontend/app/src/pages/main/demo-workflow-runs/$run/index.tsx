import { Separator } from '@/components/ui/separator';
import api, {
  StepRun,
  StepRunStatus,
  WorkflowRunStatus,
  queries,
} from '@/lib/api';
import CronPrettifier from 'cronstrue';
import { useMutation, useQuery } from '@tanstack/react-query';
import { Link, useOutletContext, useParams } from 'react-router-dom';
import invariant from 'tiny-invariant';
import { timeBetween } from '@/lib/utils';
import { Loading } from '@/components/ui/loading.tsx';
import { TenantContextType } from '@/lib/outlet';
import WorkflowRunVisualizer from './components/workflow-run-visualizer';
import { useEffect, useState } from 'react';
import { StepRunPlayground } from './components/step-run-playground';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import RelativeDate from '@/components/molecules/relative-date';
import { RunStatus } from '../components/run-statuses';
import { Button } from '@/components/ui/button';

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { BiDotsVertical } from 'react-icons/bi';
import { useApiError } from '@/lib/hooks';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';

export const WORKFLOW_RUN_TERMINAL_STATUSES = [
  WorkflowRunStatus.CANCELLED,
  WorkflowRunStatus.FAILED,
  WorkflowRunStatus.SUCCEEDED,
];

export default function ExpandedWorkflowRun() {
  const [selectedStepRun, setSelectedStepRun] = useState<StepRun | undefined>();

  const { tenant } = useOutletContext<TenantContextType>();
  invariant(tenant);

  const params = useParams();
  invariant(params.run);

  const runQuery = useQuery({
    ...queries.workflowRuns.get(tenant.metadata.id, params.run),
    refetchInterval: (query) => {
      const data = query.state.data;
      if (
        data &&
        data.jobRuns &&
        data.jobRuns.some((x) => x.status === 'RUNNING')
      ) {
        return 1000;
      }
    },
  });

  // select the first step run by default
  useEffect(() => {
    if (
      !selectedStepRun &&
      runQuery.data &&
      runQuery.data.jobRuns &&
      runQuery.data.jobRuns[0].stepRuns
    ) {
      setSelectedStepRun(runQuery.data.jobRuns[0].stepRuns[0]);
    }

    // if there is a selected step run, make sure it's still in the list
    if (
      selectedStepRun &&
      runQuery.data &&
      runQuery.data.metadata.id === params.run &&
      runQuery.data.jobRuns &&
      runQuery.data.jobRuns[0].stepRuns
    ) {
      const stepRun = runQuery.data.jobRuns.find((jobRun) =>
        jobRun.stepRuns?.find(
          (stepRun) => stepRun.metadata.id === selectedStepRun.metadata.id,
        ),
      );

      if (!stepRun) {
        setSelectedStepRun(runQuery.data.jobRuns[0].stepRuns[0]);
      }
    }
  }, [runQuery.data, params.run, selectedStepRun]);

  const { handleApiError } = useApiError({});

  const cancelWorkflowRunMutation = useMutation({
    mutationKey: [
      'workflow-run:cancel',
      runQuery?.data?.tenantId,
      runQuery?.data?.metadata.id,
    ],
    mutationFn: async () => {
      const tenantId = runQuery?.data?.tenantId;
      const workflowRunId = runQuery?.data?.metadata.id;

      invariant(tenantId, 'has tenantId');
      invariant(workflowRunId, 'has tenantId');

      const res = await api.workflowRunCancel(tenantId, {
        workflowRunIds: [workflowRunId],
      });

      return res.data;
    },
    onError: handleApiError,
  });

  if (runQuery.isLoading || !runQuery.data) {
    return <Loading />;
  }

  const run = runQuery.data;

  // const ParentLink: React.FC<{ parentId: string }> = ({ parentId }) => {
  //   return (
  //     <a
  //       href={`/workflow-runs/${parentId}`}
  //       className="flex flex-row gap-2 items-center"
  //     >
  //       <BiGitBranch />
  //       <span className="text-sm text-gray-700 dark:text-gray-300">
  //         Parent workflow
  //       </span>
  //     </a>
  //   );
  // };

  return (
    <div className="flex-grow h-full w-full">
      <div className="flex flex-col mx-auto gap-2 max-w-7xl px-4 sm:px-6 lg:px-8">
        {/* TODO the triggeredBy parent id is itself */}
        {/* {run?.triggeredBy?.parentId && <ParentLink parentId={run.triggeredBy.parentId} />} */}
        <div className="flex flex-col sm:flex-row justify-between items-center">
          <div className="flex flex-row gap-4 items-center">
            <h2 className="text-lg sm:text-2xl font-bold leading-tight text-foreground flex flex-row  items-center">
              <Link
                to={`/workflows/${run?.workflowVersion?.workflow?.metadata.id}`}
              >
                {run?.workflowVersion?.workflow?.name}-
                {run?.displayName?.split('-')[1] || run?.metadata.id}
              </Link>
              /{selectedStepRun?.step?.readableId || '*'}
            </h2>
          </div>
          <div className="flex flex-row gap-2 items-center justify-between">
            <RunStatus
              status={run.status}
              className="text-sm mt-1 px-4 shrink"
            />
            <span className="hidden sm:block">
              <DropdownMenu>
                <DropdownMenuTrigger>
                  <Button
                    aria-label="Workflow Actions"
                    size="icon"
                    variant="outline"
                  >
                    <BiDotsVertical />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent>
                  <DropdownMenuItem
                    disabled={WORKFLOW_RUN_TERMINAL_STATUSES.includes(
                      run.status,
                    )}
                    onClick={() => {
                      cancelWorkflowRunMutation.mutate();
                    }}
                  >
                    Cancel all running steps
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </span>
          </div>
        </div>
        <div className="flex flex-row justify-start items-center gap-2">
          <div className="text-sm text-gray-700 dark:text-gray-300">
            Created <RelativeDate date={run?.metadata.createdAt} />
          </div>
          {run?.startedAt && (
            <div className="text-sm text-gray-700 dark:text-gray-300">
              Started <RelativeDate date={run.startedAt} />
            </div>
          )}
          {run?.startedAt && run?.finishedAt && (
            <div className="text-sm text-gray-700 dark:text-gray-300">
              Duration {timeBetween(run.startedAt, run.finishedAt)}
            </div>
          )}
        </div>
        {run.triggeredBy?.cronSchedule && (
          <TriggeringCronSection cron={run.triggeredBy.cronSchedule} />
        )}
        <Tabs defaultValue="overview">
          <TabsList layout="underlined">
            <TabsTrigger variant="underlined" value="overview">
              Overview
            </TabsTrigger>
          </TabsList>
          <TabsContent value="overview">
            <Alert variant="default" className="my-4">
              <AlertTitle className="font-semibold">
                <span className="mr-1">🪓</span> Here's your workflow.
              </AlertTitle>
              <AlertDescription>
                Here's the same DAG, but updated with real time state of how
                your workflow is executing! You can click through each step to
                see inputs and outputs or replay the step. We've added some
                artificial delay in each step so you can see the steps in
                action.
              </AlertDescription>
            </Alert>
            <div className="w-full h-[200px] mt-8">
              <WorkflowRunVisualizer
                workflowRun={run}
                selectedStepRun={selectedStepRun}
                setSelectedStepRun={(step) => {
                  setSelectedStepRun(
                    step.stepId === selectedStepRun?.stepId ? undefined : step,
                  );
                }}
              />
            </div>
            <Separator className="my-4" />
            {!selectedStepRun ? (
              'Select a step to rerun and view details.'
            ) : (
              <StepRunPlayground
                stepRun={selectedStepRun}
                setStepRun={setSelectedStepRun}
                workflowRun={run}
                demo={true}
              />
            )}
          </TabsContent>
        </Tabs>
        <Separator className="my-4" />
      </div>
    </div>
  );
}

const getStatusText = (stepRun: StepRun): string => {
  switch (stepRun.status) {
    case StepRunStatus.RUNNING:
      return 'This step is currently running';
    case StepRunStatus.FAILED:
      return stepRun.error
        ? `This step failed with error ${stepRun.error}`
        : 'This step failed';
    case StepRunStatus.CANCELLED:
      return getCancelledStatusText(stepRun);
    case StepRunStatus.SUCCEEDED:
      return 'This step succeeded';
    case StepRunStatus.PENDING:
      return 'This step is pending';
    default:
      return 'Unknown';
  }
};

const getCancelledStatusText = (stepRun: StepRun): string => {
  switch (stepRun.cancelledReason) {
    case 'CANCELLED_BY_USER':
      return 'This step was cancelled by a user';
    case 'TIMED_OUT':
      return `This step was cancelled because it exceeded its timeout of ${stepRun.step?.timeout || '60s'}`;
    case 'SCHEDULING_TIMED_OUT':
      return `This step was cancelled because no workers were available to run ${stepRun.step?.action}`;
    case 'PREVIOUS_STEP_TIMED_OUT':
      return 'This step was cancelled because the previous step timed out';
    default:
      return `This step was cancelled (${stepRun.cancelledReason})`;
  }
};

export const StepStatusDetails = ({ stepRun }: { stepRun: StepRun }) => {
  return getStatusText(stepRun);
};

export function StepStatusSection({ stepRun }: { stepRun: StepRun }) {
  const statusText = StepStatusDetails({ stepRun });

  return (
    <div className="mb-4">
      <h3 className="font-semibold leading-tight text-foreground mb-4">
        Status
      </h3>
      <div className="text-sm text-gray-700 dark:text-gray-300">
        {statusText}
      </div>
    </div>
  );
}

export function StepDurationSection({ stepRun }: { stepRun: StepRun }) {
  return (
    <div className="mb-4">
      <h3 className="font-semibold leading-tight text-foreground mb-4">
        Duration
      </h3>
      <div className="text-sm text-gray-700 dark:text-gray-300">
        {stepRun.startedAt &&
          stepRun.finishedAt &&
          timeBetween(stepRun.startedAt, stepRun.finishedAt)}
      </div>
    </div>
  );
}

export function StepConfigurationSection({ stepRun }: { stepRun: StepRun }) {
  return (
    <div className="mb-4">
      <h3 className="font-semibold leading-tight text-foreground mb-4">
        Configuration
      </h3>
      <div className="text-sm text-gray-700 dark:text-gray-300">
        Timeout: {stepRun.step?.timeout || '60s'}
      </div>
    </div>
  );
}

function TriggeringCronSection({ cron }: { cron: string }) {
  const prettyInterval = `runs ${CronPrettifier.toString(
    cron,
  ).toLowerCase()} UTC`;

  return (
    <div className="text-sm text-gray-700 dark:text-gray-300">
      Triggered by cron {cron} which {prettyInterval}
    </div>
  );
}
