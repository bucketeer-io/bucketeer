import { useMemo } from 'react';
import { Trans } from 'react-i18next';
import { Link } from 'react-router-dom';
import { getCurrentEnvironment, useAuth } from 'auth';
import {
  PAGE_PATH_EXPERIMENTS,
  PAGE_PATH_FEATURE_AUTOOPS,
  PAGE_PATH_FEATURES
} from 'constants/routing';
import { useTranslation } from 'i18n';
import { Goal } from '@types';
import { IconExperimentsConnected, IconOperationsConnected } from '@icons';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';

export type ConfirmModalProps = {
  goal: Goal;
  isOpen: boolean;
  onClose: () => void;
};

const ConnectionsModal = ({ goal, isOpen, onClose }: ConfirmModalProps) => {
  const { t } = useTranslation(['common']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const connectionType = useMemo(() => goal?.connectionType, [goal]);
  const isExperimentType = connectionType === 'EXPERIMENT';
  const connections = isExperimentType ? goal.experiments : goal.autoOpsRules;

  return (
    <DialogModal
      className="max-w-[300px] sm:max-w-[500px]"
      title={
        isExperimentType
          ? t(`experiments-connected`)
          : t(`operations-connected`)
      }
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="flex-center flex-col w-full items-start px-5 py-8 gap-y-8">
        <div className="flex-center w-full">
          <Icon
            icon={
              isExperimentType
                ? IconExperimentsConnected
                : IconOperationsConnected
            }
            size={'fit'}
          />
        </div>

        <div className="flex-center flex-col w-full gap-y-5">
          <div className="flex-center w-full text-center px-0 sm:px-[67px] text-gray-700">
            <Trans
              i18nKey="goal-connected-desc"
              values={{
                type: t(
                  isExperimentType ? 'source-type.experiment' : 'operation'
                )
              }}
            />
          </div>
          <div className="flex flex-col w-full p-4 gap-y-3 rounded bg-gray-100 max-h-[300px] overflow-auto">
            {connections?.map((item, index) => (
              <div
                key={index}
                className="flex items-center gap-x-2 typo-para-medium text-primary-500"
              >
                <p>{index + 1}.</p>
                <Link
                  className="underline line-clamp-1 break-all"
                  to={
                    isExperimentType
                      ? `/${currentEnvironment.urlCode}${PAGE_PATH_EXPERIMENTS}/${item.id}`
                      : `/${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${item.featureId}${PAGE_PATH_FEATURE_AUTOOPS}`
                  }
                >
                  {item.featureName}
                </Link>
              </div>
            ))}
          </div>
        </div>
      </div>

      <ButtonBar
        primaryButton={
          <Button type="button" onClick={onClose}>
            {t(`close`)}
          </Button>
        }
      />
    </DialogModal>
  );
};

export default ConnectionsModal;
