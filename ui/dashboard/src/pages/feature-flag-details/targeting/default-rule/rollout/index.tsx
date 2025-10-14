import { Trans } from 'react-i18next';
import { IconMoreHorizOutlined } from 'react-icons-material-design';
import { useTranslation } from 'i18n';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import Dropdown from 'components/dropdown';
import Input from 'components/input';
import InputGroup from 'components/input-group';
import { Popover } from 'components/popover';
import PercentageInput from '../../segment-rule/percentage-input';

const DefaultRuleRollout = () => {
  const { t } = useTranslation(['form', 'table', 'common']);
  const options = [
    {
      label: 'Hours',
      value: 'hour'
    },
    {
      label: 'Daily',
      value: 'day'
    },
    {
      label: 'Weekly',
      value: 'week'
    }
  ];

  return (
    <div className="flex flex-col w-full gap-y-6">
      <p className="typo-para-medium text-gray-700">
        {t('common:source-type.progressive-rollout')}
      </p>
      <div className="flex items-center w-full gap-x-2">
        <p className="typo-para-medium text-gray-600">
          {t('table:results.variation')}
        </p>
        <FlagVariationPolygon index={0} />
        <p className="typo-para-medium text-gray-600">
          {t('table:results.variation')}
        </p>
      </div>
      <div className="flex flex-col w-full gap-y-2">
        <p className="typo-para-small text-gray-600">
          {t('targeting.rollout-to')}
        </p>
        <div className="flex items-center gap-x-3">
          <PercentageInput
            name={``}
            showVariationName={false}
            handleChangeRolloutWeight={value => {
              console.log(value);
            }}
          />
          <p className="typo-para-small text-gray-600">{t('targeting.for')}</p>
          <div className="flex items-center gap-x-2">
            <InputGroup
              addon={
                <Dropdown
                  labelCustom="4"
                  className="size-full !border-l border-r-0 border-y-0 !border-gray-400 rounded-l-none !shadow-none"
                  options={options}
                />
              }
              addonSlot="right"
              className="w-[159px] overflow-hidden"
              addonClassName="top-[1px] bottom-[1px] right-[1px] translate-x-0 translate-y-0 !flex-center rounded-r-lg w-[100px] typo-para-medium text-gray-700"
            >
              <Input
                onWheel={e => e.currentTarget.blur()}
                type="number"
                className="text-right pl-[5px] pr-[112px]"
              />
            </InputGroup>
            <Popover
              icon={IconMoreHorizOutlined}
              align="start"
              options={[
                { label: 'Add', value: 'add' },
                { label: 'Delete', value: 'delete' }
              ]}
            />
          </div>
        </div>
        <div className="flex items-center gap-x-2 pl-3 pt-2">
          <p className="typo-para-medium text-gray-600 w-16">100%</p>
          <p className="typo-para-medium text-gray-600">
            <Trans
              i18nKey={'form:targeting.total-duration'}
              components={{
                b: <span className="typo-head-bold-small text-gray-700" />
              }}
              values={{
                value: `20 hours`
              }}
            />
          </p>
        </div>
      </div>
    </div>
  );
};

export default DefaultRuleRollout;
