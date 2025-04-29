import { useCallback, useEffect, useState } from 'react';
import { getCurrentEnvironment, useAuth } from 'auth';
import { cloneDeep } from 'lodash';
import { Experiment, ExperimentResult, Feature } from '@types';
import PageLayout from 'elements/page-layout';
import GoalResultItem from './goal-results';
import { EmptyCollection } from './results-empty';

export type GoalResultTab = 'EVALUATION' | 'CONVERSION';
export type ChartDataType =
  | 'evaluation-user'
  | 'goal-total'
  | 'goal-user'
  | 'conversion-rate'
  | 'value-total'
  | 'value-user';

export interface GoalResultState {
  tab: GoalResultTab;
  chartType: ChartDataType;
  goalId: string;
}

const Results = ({
  isLoading,
  isErrorState,
  feature,
  experiment,
  experimentResult
}: {
  isLoading: boolean;
  isErrorState: boolean;
  feature?: Feature;
  experiment: Experiment;
  experimentResult?: ExperimentResult;
}) => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const [goalResultState, setGoalResultState] = useState<GoalResultState[]>([]);
  const [goalsNarrow, setGoalsNarrow] = useState<string[]>([]);

  const handleChangeResultState = ({
    index,
    tab,
    chartType
  }: {
    index: number;
    tab?: GoalResultTab;
    chartType?: ChartDataType;
  }) => {
    const cloneGoalResultState = cloneDeep(goalResultState);
    cloneGoalResultState[index] = {
      ...cloneGoalResultState[index],
      ...(tab ? { tab } : {}),
      ...(chartType ? { chartType } : {})
    };
    setGoalResultState(cloneGoalResultState);
  };

  const handleNarrowGoalResult = useCallback(
    (goalId: string) => {
      const isExisted = goalsNarrow.includes(goalId);
      setGoalsNarrow(
        isExisted
          ? goalsNarrow.filter(item => item !== goalId)
          : [...goalsNarrow, goalId]
      );
    },
    [goalsNarrow]
  );

  useEffect(() => {
    if (experimentResult?.goalResults?.length) {
      const _goalResultState = experimentResult?.goalResults.map(
        item =>
          ({
            tab: 'EVALUATION',
            chartType: 'evaluation-user',
            goalId: item.goalId
          }) as GoalResultState
      );
      setGoalResultState(_goalResultState);
    }
  }, [experimentResult]);

  return isLoading ? (
    <PageLayout.LoadingState />
  ) : isErrorState ? (
    <EmptyCollection />
  ) : (
    <div className="flex flex-col w-full gap-y-6">
      {experimentResult?.goalResults?.map((item, index) => (
        <GoalResultItem
          key={index}
          isNarrow={goalsNarrow.includes(item.goalId)}
          experiment={experiment}
          feature={feature}
          environmentId={currentEnvironment.id}
          isRequireComment={currentEnvironment.requireComment}
          goalResult={item}
          goalResultState={goalResultState[index]}
          onChangeResultState={(tab, chartType) =>
            handleChangeResultState({ index, tab, chartType })
          }
          handleNarrowGoalResult={handleNarrowGoalResult}
        />
      ))}
    </div>
  );
};

export default Results;
