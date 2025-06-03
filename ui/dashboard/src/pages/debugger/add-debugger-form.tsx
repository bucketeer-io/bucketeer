import { useTranslation } from 'i18n';
import { Evaluation, Feature } from '@types';
import Button from 'components/button';
import DebuggerAttributes from './debugger-attributes';
import DebuggerFlags from './debugger-flags';
import DebuggerUserIds from './debugger-user-ids';

const AddDebuggerForm = ({
  isOnTargeting,
  isLoading,
  evaluations,
  feature,
  onCancel
}: {
  isOnTargeting?: boolean;
  isLoading: boolean;
  evaluations: Evaluation[];
  feature?: Feature;
  onCancel?: () => void;
}) => {
  const { t } = useTranslation(['common']);
  return (
    <div className="flex flex-col w-full gap-y-6 p-6">
      <DebuggerFlags isOnTargeting={isOnTargeting} feature={feature} />
      <DebuggerUserIds />
      <DebuggerAttributes />
      {!isOnTargeting && (
        <div className="flex items-center w-full gap-x-4">
          {evaluations.length > 0 && (
            <Button
              variant={'secondary-2'}
              className="w-fit"
              onClick={onCancel}
            >
              {t('cancel')}
            </Button>
          )}
          <Button className="w-fit" loading={isLoading}>
            {t('evaluate')}
          </Button>
        </div>
      )}
    </div>
  );
};

export default AddDebuggerForm;
