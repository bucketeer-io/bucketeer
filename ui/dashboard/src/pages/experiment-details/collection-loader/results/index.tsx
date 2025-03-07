import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { useQueryExperimentResultDetails } from '@queries/experiment-result';
import { getCurrentEnvironment, useAuth } from 'auth';
import { cloneDeep } from 'lodash';
import { Experiment } from '@types';
import { EmptyCollection } from 'pages/experiment-details/page-empty';
import PageLayout from 'elements/page-layout';
import GoalResultItem from './goal-results';

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

const Results = ({ experiment }: { experiment: Experiment }) => {
  const params = useParams();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const [goalResultState, setGoalResultState] = useState<GoalResultState[]>([]);

  const {
    data: experimentResultCollection,
    isLoading,
    isError
  } = useQueryExperimentResultDetails({
    params: {
      experimentId: params?.experimentId || '',
      environmentId: currentEnvironment.id
    }
  });

  const experimentResult = experimentResultCollection?.experimentResult;
  const isErrorState = isError || !experimentResult;

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

  return (
    <>
      {isLoading ? (
        <PageLayout.LoadingState />
      ) : isErrorState ? (
        <EmptyCollection />
      ) : (
        <div className="flex flex-col w-full gap-y-6">
          {experimentResult?.goalResults?.map((item, index) => (
            <GoalResultItem
              key={index}
              experiment={experiment}
              goalResult={item}
              goalResultState={goalResultState[index]}
              onChangeResultState={(tab, chartType) =>
                handleChangeResultState({ index, tab, chartType })
              }
            />
          ))}
        </div>
      )}
    </>
  );
};

export default Results;
