import { Experiment } from '@types';
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
}) => (
  <div className="flex flex-col size-full gap-y-6">
    <ExperimentState experiment={experiment} />
    {currentTab === 'results' && <Results experiment={experiment} />}
    {currentTab === 'settings' && (
      <ExperimentSettings experiment={experiment} />
    )}
  </div>
);

export default CollectionLoader;
