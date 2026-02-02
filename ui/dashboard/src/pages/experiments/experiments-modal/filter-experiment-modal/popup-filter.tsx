import { FilterTypes } from 'hooks/use-options';
import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import { IconPlus, IconTrash } from '@icons';
import { ExperimentFilters } from 'pages/experiments/types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
import Dropdown, { DropdownOption, DropdownValue } from 'components/dropdown';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';
import useExperimentFilterLogic from './use-filter-experiment-logic';

export type FilterProps = {
  onSubmit: (v: Partial<ExperimentFilters>) => void;
  isOpen: boolean;
  onClose: () => void;
  onClearFilters: () => void;
  filters?: Partial<ExperimentFilters>;
};

const FilterExperimentPopup = ({
  onSubmit,
  isOpen,
  onClose,
  onClearFilters,
  filters
}: FilterProps) => {
  const { t } = useTranslation(['common']);

  const {
    selectedFilters,
    isDisabledAddButton,
    isDisabledSubmitButton,
    isLoading,
    remainingFilterOptions,
    setSelectedFilters,
    onConfirmHandler,
    getValueOptions,
    handleGetLabelFilterValue,
    handleChangeFilterValue,
    handleChangeOption
  } = useExperimentFilterLogic(filters, onSubmit);

  return (
    <DialogModal
      className="w-[750px]"
      title={t('filters')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="flex flex-col w-full items-start p-5 gap-y-4">
        {selectedFilters.map((filterOption, filterIndex) => {
          const { label, value: filterType } = filterOption;
          const isMaintainerFilter = filterType === FilterTypes.MAINTAINER;
          const isStatusFilter = filterType === FilterTypes.STATUSES;
          const valueOptions = getValueOptions(filterOption);

          return (
            <div
              key={filterIndex}
              className="flex items-center w-full h-12 gap-x-4"
            >
              <div
                className={cn(
                  'typo-para-small text-center py-[3px] w-[42px] min-w-[42px] rounded text-accent-pink-500 bg-accent-pink-50',
                  {
                    'bg-gray-200 text-gray-600': filterIndex !== 0
                  }
                )}
              >
                {t(filterIndex === 0 ? `if` : 'and')}
              </div>
              <Divider vertical={true} className="border-primary-500" />
              <Dropdown
                placeholder={t(`select-filter`)}
                labelCustom={label}
                className="w-full"
                contentClassName="w-[235px]"
                options={remainingFilterOptions as DropdownOption[]}
                value={filterType}
                onChange={value =>
                  handleChangeOption(value as DropdownValue, filterIndex)
                }
              />

              <p className="typo-para-medium text-gray-600">is</p>
              <Dropdown
                disabled={isLoading || !filterType}
                placeholder={t(`select-value`)}
                labelCustom={handleGetLabelFilterValue(filterOption)}
                className={cn('w-full max-w-[280px] truncate', {
                  capitalize: isStatusFilter
                })}
                contentClassName={cn('w-[235px]', {
                  'pt-0 w-[300px]': isMaintainerFilter,
                  'hidden-scroll': valueOptions?.length > 15
                })}
                value={
                  isStatusFilter && Array.isArray(filterOption?.filterValue)
                    ? filterOption.filterValue
                    : (filterOption.filterValue as string)
                }
                options={valueOptions as DropdownOption[]}
                isSearchable={isMaintainerFilter}
                onChange={value =>
                  handleChangeFilterValue(value as string, filterIndex)
                }
                multiselect={isStatusFilter}
                isListItem={isMaintainerFilter || isStatusFilter}
              />

              <Button
                variant={'grey'}
                className="px-0 w-fit"
                disabled={selectedFilters.length <= 1}
                onClick={() =>
                  setSelectedFilters(
                    selectedFilters.filter((_, index) => filterIndex !== index)
                  )
                }
              >
                <Icon icon={IconTrash} size={'sm'} />
              </Button>
            </div>
          );
        })}
        <Button
          disabled={isDisabledAddButton}
          variant={'text'}
          className="h-6 px-0"
          onClick={() => {
            setSelectedFilters([
              ...selectedFilters,
              {
                label: '',
                value: undefined,
                filterValue: ''
              }
            ]);
          }}
        >
          <Icon icon={IconPlus} />
          {t('add-filter')}
        </Button>
      </div>

      <ButtonBar
        secondaryButton={
          <Button disabled={isDisabledSubmitButton} onClick={onConfirmHandler}>
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

export default FilterExperimentPopup;
