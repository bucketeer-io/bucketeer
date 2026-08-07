import { Trans } from 'react-i18next';
import { Link } from 'react-router';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_FEATURES } from 'constants/routing';
import { useTranslation } from 'i18n';
import { Feature } from '@types';
import { IconToastWarning } from '@icons';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';

export type ConfirmRulesImpactModalProps = {
  features: Feature[];
  isOpen: boolean;
  loading?: boolean;
  onClose: () => void;
  onSubmit: () => void;
};

// Shown when saving a segment whose rules changed while the segment is
// referenced by feature flags: changing the rules changes who is in the
// segment, which affects targeting in every connected flag.
const ConfirmRulesImpactModal = ({
  features,
  isOpen,
  loading,
  onClose,
  onSubmit
}: ConfirmRulesImpactModalProps) => {
  const { t } = useTranslation(['common', 'form']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  return (
    <DialogModal
      className="w-[500px]"
      title={t('form:segment-rules.impact-title')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="flex flex-col w-full px-5 py-8 gap-y-5">
        <div className="flex flex-col w-full px-4 py-3 bg-accent-yellow-50 border-l-4 border-accent-yellow-500 rounded">
          <div className="flex gap-x-2 w-full pr-3">
            <Icon
              icon={IconToastWarning}
              size={'xxs'}
              color="accent-yellow-500"
              className="mt-1"
            />
            <Trans
              i18nKey={'form:segment-rules.impact-description'}
              values={{ count: features.length }}
              components={{
                p: <p className="typo-para-medium text-accent-yellow-500" />
              }}
            />
          </div>
          <div className="flex flex-col w-full gap-y-1 mt-2">
            {features?.map((item, index) => (
              <div
                key={item.id}
                className="flex gap-x-2 w-full pl-6 typo-para-medium text-primary-500"
              >
                <p>{index + 1}.</p>
                <Link
                  target="_blank"
                  to={`/${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${item.id}/targeting`}
                  className="hover:underline truncate"
                >
                  {item.name}
                </Link>
              </div>
            ))}
          </div>
        </div>
      </div>
      <ButtonBar
        secondaryButton={
          <Button loading={loading} className="w-24" onClick={onSubmit}>
            {t(`submit`)}
          </Button>
        }
        primaryButton={
          <Button onClick={onClose} variant="secondary">
            {t(`cancel`)}
          </Button>
        }
      />
    </DialogModal>
  );
};

export default ConfirmRulesImpactModal;
