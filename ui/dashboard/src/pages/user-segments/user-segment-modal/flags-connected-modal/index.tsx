import { Link } from 'react-router-dom';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_FEATURES } from 'constants/routing';
import { useTranslation } from 'i18n';
import { UserSegment } from '@types';
import { IconFlagConnected } from '@icons';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';

export type DeleteMemberProps = {
  segment: UserSegment;
  isOpen: boolean;
  onClose: () => void;
};

const FlagsConnectedModal = ({
  segment,
  isOpen,
  onClose
}: DeleteMemberProps) => {
  const { t } = useTranslation(['common']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  return (
    <DialogModal
      className="w-[496px]"
      title={t(`flags-connected`)}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="flex flex-col w-full px-5 py-8 gap-y-5">
        <div className="flex-center flex-col gap-8">
          <Icon icon={IconFlagConnected} size={'fit'} />
          <div className="typo-para-medium text-gray-700 text-center px-[42px]">
            {t('flags-connected-desc')}
          </div>
        </div>
        <div className="flex flex-col w-full gap-y-5 p-4 bg-primary-50 rounded">
          {segment?.features?.map((item, index) => (
            <div
              key={index}
              className="flex items-center gap-x-2 w-full truncate typo-para-medium text-primary-500"
            >
              <p>{index + 1}.</p>
              <Link
                to={`/${currentEnvironment.urlCode}/${PAGE_PATH_FEATURES}/${item.id}/targeting`}
                className="underline"
              >
                {item.name}
              </Link>
            </div>
          ))}
        </div>
      </div>
      <ButtonBar
        primaryButton={
          <Button onClick={onClose} variant="primary">
            {t(`close`)}
          </Button>
        }
      />
    </DialogModal>
  );
};

export default FlagsConnectedModal;
