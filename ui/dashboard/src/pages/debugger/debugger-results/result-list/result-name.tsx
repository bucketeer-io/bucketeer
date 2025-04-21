import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_FEATURES, PAGE_PATH_MEMBERS } from 'constants/routing';
import { Feature, FeatureVariationType } from '@types';
import { IconUserOutlined } from '@icons';
import {
  FlagDataTypeIcon,
  FlagNameElement
} from 'pages/feature-flags/collection-layout/elements';
import {
  getDataTypeIcon,
  getFlagStatus
} from 'pages/feature-flags/collection-layout/elements/utils';

interface Props {
  isFlag: boolean;
  id: string;
  name: string;
  variationType: FeatureVariationType;
  feature: Feature;
  maintainer: string;
  onTable?: boolean;
}

const ResultName = ({
  isFlag,
  id,
  name,
  variationType,
  feature,
  maintainer,
  onTable
}: Props) => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  return (
    <FlagNameElement
      className={onTable ? 'col-span-4' : ''}
      id={id}
      name={name}
      icon={getDataTypeIcon(variationType)}
      status={isFlag ? getFlagStatus(feature) : undefined}
      iconElement={
        isFlag ? undefined : <FlagDataTypeIcon icon={IconUserOutlined} />
      }
      link={
        isFlag
          ? `/${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${id}/targeting`
          : `/${currentEnvironment.urlCode}${PAGE_PATH_MEMBERS}/`
      }
      variationType={variationType}
      maintainer={maintainer}
    />
  );
};

export default ResultName;
