import { Trans } from 'react-i18next';
import {
  IconArrowDownwardFilled,
  IconArrowUpwardFilled
} from 'react-icons-material-design';
import { useTranslation } from 'i18n';
import { IconInfo } from '@icons';
import Icon from 'components/icon';
import Card from '../../elements/card';
import AddRuleButton from '../add-rule-button';
import Condition from './condition';

const TargetSegmentRule = () => {
  const { t } = useTranslation(['table']);

  return (
    <Card>
      <div>
        <div className="flex items-center gap-x-2">
          <p className="typo-para-medium leading-4 text-gray-700">
            {t('feature-flags.rules')}
          </p>
          <Icon icon={IconInfo} size={'xxs'} color="gray-500" />
        </div>
      </div>
      <Card className="shadow-none border border-gray-400">
        <div className="flex items-center justify-between w-full">
          <p className="typo-para-medium leading-5 text-gray-700">
            <Trans
              i18nKey={'table:feature-flags.rule-index'}
              values={{
                index: '1'
              }}
            />
          </p>
          <div className="flex items-center gap-x-1">
            <Icon icon={IconArrowDownwardFilled} color="gray-500" size={'sm'} />
            <Icon icon={IconArrowUpwardFilled} color="gray-500" size={'sm'} />
          </div>
        </div>
        <Condition type="if" situation="compare" />
        <AddRuleButton />
      </Card>
    </Card>
  );
};

export default TargetSegmentRule;
