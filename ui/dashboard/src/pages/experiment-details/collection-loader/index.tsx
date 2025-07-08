import { useParams } from 'react-router-dom';
import { useQueryExperimentResultDetails } from '@queries/experiment-result';
import { useQueryFeature } from '@queries/feature-details';
import { getCurrentEnvironment, useAuth } from 'auth';
import { Experiment } from '@types';
import { checkEnvironmentEmptyId } from 'utils/function';
import { ExperimentDetailsTab } from '../page-content';
import Results from './results';
import ExperimentSettings from './settings';
import ExperimentState from './settings/experiment-state';

const CollectionLoader = ({
  currentTab,
  experiment
}: {
  currentTab: ExperimentDetailsTab;
  experiment: Experiment;
}) => {
  const params = useParams();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const {
    data: experimentResultCollection,
    isLoading,
    isError
  } = useQueryExperimentResultDetails({
    params: {
      experimentId: params?.experimentId || '',
      environmentId: checkEnvironmentEmptyId(currentEnvironment.id)
    },
    retry: experiment.status !== 'WAITING'
  });

  const experimentResult = experimentResultCollection?.experimentResult;

  const { data: featureResultCollection, isError: featureError } =
    useQueryFeature({
      params: {
        environmentId: checkEnvironmentEmptyId(currentEnvironment.id),
        id: experiment.featureId
      }
    });
  const feature = featureResultCollection?.feature;

  const isErrorState = isError || !experimentResult || featureError;

  return (
    <div className="flex flex-col size-full gap-y-6">
      <ExperimentState
        experiment={experiment}
        experimentResult={experimentResult}
      />
      {currentTab === 'results' && (
        <Results
          isLoading={isLoading}
          isErrorState={isErrorState}
          experiment={experiment}
          experimentResult={experimentResult}
          feature={feature}
        />
      )}
      {currentTab === 'settings' && (
        <ExperimentSettings experiment={experiment} />
      )}
    </div>
  );
};

export default CollectionLoader;
