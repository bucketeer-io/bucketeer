import { useTranslation } from 'i18n';
import { UserSegmentsFilters } from 'pages/user-segments/types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
import Dropdown from 'components/dropdown';
import SlideModal from 'components/modal/slide';
import useFilterSegmentLogic from './use-filter-segment-logic';

export type FilterProps = {
  onSubmit: (v: Partial<UserSegmentsFilters>) => void;
  isOpen: boolean;
  onClose: () => void;
  onClearFilters: () => void;
  filters?: Partial<UserSegmentsFilters>;
};

const FilterUserSegmentSlideModal = ({
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
    <SlideModal title={t('filters')} isOpen={isOpen} onClose={onClose}>
      <div className="w-full h-full flex flex-col justify-between relative">
        <div className="flex flex-col w-full items-start p-5 gap-y-4">
          <div className="flex items-start w-full h-[100px] gap-x-3">
            <div className="h-full flex flex-col gap-y-4 items-center justify-center">
              <div className="mt-2 typo-para-small text-center py-[3px] w-[42px] min-w-[42px] rounded text-accent-pink-500 bg-accent-pink-50">
                {t(`if`)}
              </div>
              <Divider vertical={true} className="border-primary-500" />
            </div>
            <div className="flex flex-col w-full">
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
                className="w-full py-2"
                contentClassName="w-[270px]"
              />
              <div className="flex items-center gap-3 mt-3 pl-3">
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
                  className="w-full py-2"
                  contentClassName="w-[235px]"
                />
              </div>
            </div>
          </div>
        </div>

        <ButtonBar
          className="sticky bottom-0 left-0 bg-white"
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
      </div>
    </SlideModal>
  );
};

export default FilterUserSegmentSlideModal;
