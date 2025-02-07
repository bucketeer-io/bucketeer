import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import { IconPlus } from '@icons';
import Button from 'components/button';
import Icon from 'components/icon';

interface Props {
  isCenter?: boolean;
}

const AddRuleButton = ({ isCenter }: Props) => {
  const { t } = useTranslation(['table']);
  return (
    <div
      className={cn('flex items-center w-fit', {
        'w-full py-3 justify-center border border-dashed border-gray-200 rounded-lg':
          isCenter
      })}
    >
      <Button variant={'text'} className="gap-x-2 h-6 !p-0">
        <Icon
          icon={IconPlus}
          color="primary-500"
          className="flex-center"
          size={'sm'}
        />{' '}
        {t('feature-flags.add-rule')}
      </Button>
    </div>
  );
};

export default AddRuleButton;
