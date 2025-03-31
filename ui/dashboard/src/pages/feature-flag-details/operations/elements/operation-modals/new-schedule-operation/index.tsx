import { useFormContext } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { useTranslation } from 'i18n';
import { isTimestampArraySorted } from 'utils/function';
import { OperationForm } from 'pages/feature-flag-details/operations/form-schema';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
import SlideModal from 'components/modal/slide';
import ScheduleList from './schedule-list';

export interface OperationModalProps {
  isOpen: boolean;
  isEnabledFlag: boolean;
  onClose: () => void;
}

const NewScheduleOperationModal = ({
  isOpen,
  isEnabledFlag,
  onClose
}: OperationModalProps) => {
  const { t } = useTranslation(['form', 'common']);
  const {
    formState: { isValid, isSubmitting },
    watch,
    handleSubmit
  } = useFormContext<OperationForm>();

  const handleOnSubmit = () => {};

  const watchDatetimeClausesList = watch('datetimeClausesList');

  const isDateSorted = isTimestampArraySorted(
    watchDatetimeClausesList?.map(item => item.time?.getTime() ?? 0) || []
  );

  return (
    <SlideModal
      title={t('common:new-operation')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="flex flex-col gap-y-5 w-full p-5 pb-28">
        <div className="flex items-center gap-x-4 typo-head-bold-small text-gray-800">
          <Trans
            i18nKey={'form:feature-flags.current-state'}
            values={{
              state: t(`form:experiments.${isEnabledFlag ? 'on' : 'off'}`)
            }}
            components={{
              comp: (
                <div className="flex-center typo-para-small text-gray-600 px-2 py-[1px] border border-gray-400 rounded mb-[-4px]" />
              )
            }}
          />
        </div>
        <Divider />
        <ScheduleList isDateSorted={isDateSorted} />
      </div>
      <div className="absolute left-0 bottom-0 bg-gray-50 w-full rounded-b-lg">
        <ButtonBar
          primaryButton={
            <Button variant="secondary" onClick={onClose}>
              {t(`common:cancel`)}
            </Button>
          }
          secondaryButton={
            <Button
              type="submit"
              loading={isSubmitting}
              disabled={isValid || isSubmitting}
              onClick={handleSubmit(handleOnSubmit)}
            >
              {t(`feature-flags.create-operation`)}
            </Button>
          }
        />
      </div>
    </SlideModal>
  );
};

export default NewScheduleOperationModal;
