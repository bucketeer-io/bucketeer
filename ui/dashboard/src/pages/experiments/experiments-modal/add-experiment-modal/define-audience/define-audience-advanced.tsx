import { useTranslation } from 'react-i18next';
import { useToggleOpen } from 'hooks';
import { cn } from 'utils/style';
import { IconChevronRight } from '@icons';
import Checkbox from 'components/checkbox';
import Icon from 'components/icon';
import { DefineAudienceProps } from '.';

const DefineAudienceAdvanced = ({ field }: DefineAudienceProps) => {
  const { t } = useTranslation(['form', 'common']);

  const [isOpenAdvanced, onOpenAdvanced, onCloseAdvanced] =
    useToggleOpen(false);

  return (
    <div className="flex flex-col w-full gap-y-5">
      <div
        className="flex items-center gap-x-2 cursor-pointer"
        onClick={isOpenAdvanced ? onCloseAdvanced : onOpenAdvanced}
      >
        <p className="typo-para-medium text-gray-600 leading-5">
          {t('advanced')}
        </p>
        <Icon
          icon={IconChevronRight}
          color="gray-600"
          className={cn('transition-all duration-200 rotate-90', {
            '-rotate-90': isOpenAdvanced
          })}
        />
      </div>
      {isOpenAdvanced && (
        <div>
          <div className="flex items-center gap-x-2">
            <Checkbox
              id="variationReassignment"
              checked={field.value?.variationReassignment}
              onCheckedChange={value =>
                field.onChange({
                  ...field.value,
                  variationReassignment: value
                })
              }
            />
            <label
              htmlFor="variationReassignment"
              className="typo-para-medium text-gray-600 leading-5 cursor-pointer"
            >
              {t('experiments.define-audience.advanced-prevent-variation')}
            </label>
          </div>
          <p className="typo-para-medium text-gray-500 leading-5 mt-3">
            {t('experiments.define-audience.advanced-context-remain')}
          </p>
        </div>
      )}
    </div>
  );
};

export default DefineAudienceAdvanced;
