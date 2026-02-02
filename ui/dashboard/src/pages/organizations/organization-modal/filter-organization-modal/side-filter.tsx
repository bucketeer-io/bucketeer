import { useTranslation } from 'react-i18next';
import { OrganizationFilters } from 'pages/organizations/types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
import Dropdown, { DropdownOption } from 'components/dropdown';
import SlideModal from 'components/modal/slide';
import useFilterOrganizationLogic from './use-filter-organization-logic';

export type FilterProps = {
  isOpen: boolean;
  filters?: Partial<OrganizationFilters>;
  onSubmit: (v: Partial<OrganizationFilters>) => void;
  onClose: () => void;
  onClearFilters: () => void;
};

const FilterOrganizationSlide = ({
  isOpen,
  filters,
  onSubmit,
  onClose,
  onClearFilters
}: FilterProps) => {
  const { t } = useTranslation(['common']);
  const {
    selectedFilterType,
    filterEnabledOptions,
    isDisabledSubmitBtn,
    enabledOptions,
    selectedValue,
    setSelectedFilterType,
    setSelectedValue,
    onConfirmHandler
  } = useFilterOrganizationLogic(onSubmit, filters);
  return (
    <SlideModal title={t('filters')} isOpen={isOpen} onClose={onClose}>
      <div className="w-full h-full flex flex-col justify-between relative">
        <div className="flex flex-col w-full items-start p-5 gap-y-4">
          <div className="flex items-start w-full h-[100px] gap-x-3">
            <div className="h-full flex flex-col gap-y-4 items-center justify-center">
              <div className="mt-2 typo-para-small text-center py-[3px] w-[42px] min-w-[42px] rounded text-accent-pink-500 bg-accent-pink-50">
                {t('if')}
              </div>
              <Divider vertical={true} className="border-primary-500" />
            </div>
            <div className="flex flex-col w-full">
              <Dropdown
                placeholder={t('select-filter')}
                labelCustom={selectedFilterType?.label}
                options={filterEnabledOptions as DropdownOption[]}
                value={selectedFilterType?.value}
                onChange={value => {
                  const selected = filterEnabledOptions.find(
                    item => item.value === value
                  );
                  setSelectedFilterType(selected);
                  setSelectedValue(undefined);
                }}
                className="w-full truncate py-2"
                contentClassName="w-[270px]"
              />
              <div className="flex items-center gap-3 mt-3 pl-3">
                <p className="typo-para-medium text-gray-600">is</p>
                <Dropdown
                  disabled={!selectedFilterType}
                  placeholder={t('select-value')}
                  value={selectedValue?.value}
                  labelCustom={selectedValue?.label}
                  options={enabledOptions as DropdownOption[]}
                  onChange={value => {
                    const selected = enabledOptions.find(
                      item => item.value === value
                    );
                    setSelectedValue(selected);
                  }}
                  className="w-full truncate py-2"
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
              {t('confirm')}
            </Button>
          }
          primaryButton={
            <Button onClick={onClearFilters} variant="secondary">
              {t('clear')}
            </Button>
          }
        />
      </div>
    </SlideModal>
  );
};

export default FilterOrganizationSlide;
