import { Link } from 'react-router-dom';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_GOALS } from 'constants/routing';
import { useTranslation } from 'i18n';
import { Experiment } from '@types';
import { IconGoal } from '@icons';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';

export type ConfirmModalProps = {
  experiment: Experiment;
  isOpen: boolean;
  onClose: () => void;
};

const GoalsConnectionModal = ({
  experiment,
  isOpen,
  onClose
}: ConfirmModalProps) => {
  const { t } = useTranslation(['common']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const goals = experiment?.goalIds?.map(item =>
    experiment?.goals?.find(goal => goal.id === item)
  );

  return (
    <DialogModal
      className="w-[500px]"
      title={t(`goals-connected`)}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="flex-center flex-col w-full items-start px-5 py-8 gap-y-8">
        <div className="flex-center w-full">
          <Icon icon={IconGoal} size={'fit'} />
        </div>

        <div className="flex-center flex-col w-full gap-y-5">
          <div className="flex-center w-full text-center px-[67px] text-gray-700">
            <p className="typo-para-big">{t('goals-connected-desc')}</p>
          </div>
          <div className="flex flex-col w-full p-4 gap-y-5 rounded bg-gray-100 max-h-[300px] overflow-auto">
            {goals?.map((item, index) => (
              <div
                key={index}
                className="flex items-center gap-x-2 typo-para-medium leading-4 text-primary-500"
              >
                <p>{index + 1}.</p>
                <Link
                  className="underline line-clamp-1 break-all"
                  to={`/${currentEnvironment.urlCode}${PAGE_PATH_GOALS}/${item?.id}`}
                >
                  {item?.name}
                </Link>
              </div>
            ))}
          </div>
        </div>
      </div>

      <ButtonBar
        primaryButton={<Button onClick={onClose}>{t(`close`)}</Button>}
      />
    </DialogModal>
  );
};

export default GoalsConnectionModal;
