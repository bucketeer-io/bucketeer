import { useMemo } from 'react';
import { useFormContext } from 'react-hook-form';
import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import { IconInfo } from '@icons';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import ServeDropdown from '../serve-dropdown';

const SegmentVariation = ({
  segmentIndex,
  ruleIndex
}: {
  segmentIndex: number;
  ruleIndex: number;
}) => {
  const { t } = useTranslation(['table', 'common']);

  const methods = useFormContext();
  const { control } = methods;

  const serveOptions = useMemo(
    () => [
      { label: t('false'), value: 0 },
      { label: t('true'), value: 1 }
    ],
    []
  );

  return (
    <Form.Field
      control={control}
      name={`targetSegmentRules.${segmentIndex}.rules.${ruleIndex}.variation`}
      render={({ field }) => (
        <Form.Item className="flex flex-col flex-1 py-0 min-w-[170px]">
          <Form.Label required className="relative w-fit mb-5">
            {t('feature-flags.variation')}
            <Icon
              icon={IconInfo}
              size="xs"
              color="gray-500"
              className="absolute -right-6"
            />
          </Form.Label>
          <Form.Control>
            <ServeDropdown
              isExpand
              serveValue={field.value}
              onChangeServe={field.onChange}
            />
            {/* <div className={cn('flex items-end gap-x-6 w-full')}>
              <p className="typo-para-small text-gray-600 py-[14px] uppercase">
                {t('feature-flags.serve')}
              </p>
              <DropdownMenu>
                <div className={cn('flex flex-col gap-y-2 w-full')}>
                  <p className="typo-para-small leading-[14px] text-gray-600">
                    {t('feature-flags.variation')}
                  </p>
                  <DropdownMenuTrigger
                    label="test"
                    trigger={
                      <div className={cn('flex items-center gap-x-2')}>
                        <FlagVariationPolygon
                          color={field.value === 0 ? 'pink' : 'blue'}
                        />
                        <p className="typo-para-medium leading-5 text-gray-700">
                          {field.value === 0 ? 'False' : 'True'}
                        </p>
                      </div>
                    }
                    className={'w-full'}
                  />
                </div>
                <DropdownMenuContent align="start" {...field}>
                  {serveOptions.map((item, index) => (
                    <DropdownMenuItem
                      key={index}
                      label={item.label}
                      value={item.value}
                      icon={() => (
                        <FlagVariationPolygon
                          color={index === 0 ? 'pink' : 'blue'}
                        />
                      )}
                      onSelectOption={value => field.onChange(!!value)}
                    />
                  ))}
                </DropdownMenuContent>
              </DropdownMenu>
            </div> */}
          </Form.Control>
          <Form.Message />
        </Form.Item>
      )}
    />
  );
};

export default SegmentVariation;
