import { useMemo } from 'react';
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
  const { urlCode } = currentEnvironment;

  const isShowFlagElement = useMemo(
    () => (isFlag && !onTable) || (!isFlag && onTable),
    [isFlag, onTable]
  );

  const status = useMemo(
    () => (isShowFlagElement ? getFlagStatus(feature) : undefined),
    [isShowFlagElement, feature]
  );

  const iconElement = useMemo(
    () =>
      isShowFlagElement ? undefined : (
        <FlagDataTypeIcon icon={IconUserOutlined} />
      ),
    [isShowFlagElement]
  );

  const _maintainer = useMemo(
    () => (isShowFlagElement ? maintainer : undefined),
    [isShowFlagElement, maintainer]
  );

  const link = useMemo(
    () =>
      isShowFlagElement
        ? `/${urlCode}${PAGE_PATH_FEATURES}/${id}/targeting`
        : `/${urlCode}${PAGE_PATH_MEMBERS}/`,
    [isShowFlagElement, id, urlCode]
  );

  return (
    <FlagNameElement
      id={id}
      name={name}
      icon={getDataTypeIcon(variationType)}
      status={status}
      iconElement={iconElement}
      link={link}
      variationType={variationType}
      maintainer={_maintainer}
      maxLines={1}
      variationCls={'size-8'}
      variant="secondary"
    />
  );
};

export default ResultName;
