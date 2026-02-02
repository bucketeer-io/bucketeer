import { FilterTypes } from 'hooks/use-options';
import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import { IconPlus, IconTrash } from '@icons';
import { MembersFilters } from 'pages/members/types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
import Dropdown, { DropdownOption } from 'components/dropdown';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';
import useFilterMemberLogic from './use-filter-member-logic';

export type FilterProps = {
  onSubmit: (v: Partial<MembersFilters>) => void;
  isOpen: boolean;
  onClose: () => void;
  onClearFilters: () => void;
  filters?: Partial<MembersFilters>;
};

const FilterMemberPopup = ({
  isOpen,
  filters,
  onSubmit,
  onClose,
  onClearFilters
}: FilterProps) => {
  const { t } = useTranslation(['common']);
  const {
    selectedFilters,
    getValueOptions,
    remainingFilterOptions,
    setSelectedFilters,
    handleChangeFilterValue,
    isDisabledAddButton,
    isDisabledSubmitButton,
    onConfirmHandler,
    handleGetLabelFilterValue,
    isLoadingTeams
  } = useFilterMemberLogic(onSubmit, filters);
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
          const isTeamsFilter = filterType === FilterTypes.TEAMS;
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
                options={remainingFilterOptions as DropdownOption[]}
                value={filterType}
                onChange={value => {
                  const selected = remainingFilterOptions.find(
                    item => item.value === value
                  );
                  if (selected) {
                    selectedFilters[filterIndex] = {
                      ...selected,
                      filterValue: value === FilterTypes.TEAMS ? [] : ''
                    };
                    setSelectedFilters([...selectedFilters]);
                  }
                }}
              />

              <p className="typo-para-medium text-gray-600">is</p>
              <Dropdown
                placeholder={t(`select-value`)}
                labelCustom={handleGetLabelFilterValue(filterOption)}
                disabled={(isTeamsFilter && isLoadingTeams) || !filterType}
                loading={isTeamsFilter && isLoadingTeams}
                multiselect={isTeamsFilter}
                isSearchable={isTeamsFilter}
                value={
                  isTeamsFilter && Array.isArray(filterOption?.filterValue)
                    ? (filterOption.filterValue as string[])
                    : (filterOption.filterValue as string)
                }
                options={valueOptions as DropdownOption[]}
                onChange={value =>
                  handleChangeFilterValue(value as string, filterIndex)
                }
                className="w-full truncate"
                contentClassName={cn('w-[235px]', {
                  'pt-0 w-[300px]': isTeamsFilter,
                  'hidden-scroll': valueOptions?.length > 15
                })}
                menuContentSide="bottom"
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
                value: '',
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

export default FilterMemberPopup;
