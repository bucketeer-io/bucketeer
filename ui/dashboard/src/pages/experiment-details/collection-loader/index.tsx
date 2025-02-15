import { Experiment } from '@types';
import { ExperimentDetailsTab } from '../page-content';
import Metrics from './metrics';
import ExperimentSettings from './settings';

const CollectionLoader = ({
  currentTab,
  experiment
}: {
  currentTab: ExperimentDetailsTab;
  experiment: Experiment;
}) => {
  if (currentTab === 'results') return <Metrics />;
  if (currentTab === 'settings')
    return <ExperimentSettings experiment={experiment} />;
};

export default CollectionLoader;
