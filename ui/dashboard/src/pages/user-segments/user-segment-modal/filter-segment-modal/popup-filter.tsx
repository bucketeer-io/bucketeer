import { useTranslation } from 'i18n';
import { UserSegmentsFilters } from 'pages/user-segments/types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
import Dropdown from 'components/dropdown';
import DialogModal from 'components/modal/dialog';
import useFilterSegmentLogic from './use-filter-segment-logic';

export type FilterProps = {
  onSubmit: (v: Partial<UserSegmentsFilters>) => void;
  isOpen: boolean;
  onClose: () => void;
  onClearFilters: () => void;
  filters?: Partial<UserSegmentsFilters>;
};

const FilterUserSegmentPopup = ({
  onSubmit,
  isOpen,
  onClose,
  onClearFilters,
  filters
}: FilterProps) => {
  const { t } = useTranslation(['common']);

  const {
    selectedFilterType,
    valueOption,
    filterStatusOptions,
    segmentStatusOptions,
    isDisabledSubmitBtn,
    setSelectedFilterType,
    setValueOption,

    onConfirmHandler
  } = useFilterSegmentLogic(filters, onSubmit);

  return (
    <DialogModal
      className="max-w-[550px] lg:max-w-[750px]"
      title={t('filters')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="flex flex-col w-full items-start p-5 gap-y-4">
        <div className="flex items-center w-full h-12 gap-x-4">
          <div className="typo-para-small text-center py-[3px] px-4 rounded text-accent-pink-500 bg-accent-pink-50">
            {t(`if`)}
          </div>
          <Divider vertical={true} className="border-primary-500" />
          <Dropdown
            value={selectedFilterType?.value}
            onChange={value => {
              const selected = filterStatusOptions.find(
                item => item.value === value
              );
              setSelectedFilterType(selected);
            }}
            placeholder={t(`select-filter`)}
            options={filterStatusOptions.map(item => ({
              ...item,
              label: item.label,
              value: item.value || ''
            }))}
            className="w-full"
            contentClassName="w-[235px]"
          />

          <p className="typo-para-medium text-gray-600">is</p>
          <Dropdown
            placeholder={t(`select-value`)}
            disabled={!selectedFilterType}
            options={segmentStatusOptions.map(item => ({
              ...item,
              label: item.label,
              value: item.value || ''
            }))}
            value={valueOption?.value}
            onChange={value => {
              const selected = segmentStatusOptions.find(
                item => item.value === value
              );
              setValueOption(selected);
            }}
            className="w-full"
            contentClassName="w-[235px]"
          />
        </div>
      </div>

      <ButtonBar
        secondaryButton={
          <Button disabled={isDisabledSubmitBtn} onClick={onConfirmHandler}>
            {t(`confirm`)}
          </Button>
        }
        primaryButton={
          <Button onClick={onClearFilters} variant="secondary">
            {t(`clear`)}
          </Button>
        }
      />
    </DialogModal>
  );
};

export default FilterUserSegmentPopup;
