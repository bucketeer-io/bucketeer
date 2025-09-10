import { ReactNode } from 'react';
import { useFormContext } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { IconUndoOutlined } from 'react-icons-material-design';
import { MultiValue } from 'react-select';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import { truncateBySide } from 'utils/converts';
import { copyToClipBoard } from 'utils/function';
import { IconCopy, IconInfo } from '@icons';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import Button from 'components/button';
import { CreatableSelect, Option } from 'components/creatable-select';
import Form from 'components/form';
import Icon from 'components/icon';
import { Tooltip } from 'components/tooltip';
import Card from '../../elements/card';
import { TargetingSchema } from '../form-schema';
import { DiscardChangesType, IndividualRuleItem } from '../types';
import { getAlreadyTargetedVariation } from '../utils';

interface Props {
  individualRules: IndividualRuleItem[];
  handleDiscardChanges: (type: DiscardChangesType) => void;
}

export const UserMessage = ({ message }: { message: ReactNode }) => {
  return (
    <div className={'text-center typo-para-small text-gray-500'}>{message}</div>
  );
};

const IndividualRule = ({ individualRules, handleDiscardChanges }: Props) => {
  const { t } = useTranslation(['table', 'form', 'common', 'message']);
  const { notify } = useToast();

  const methods = useFormContext<TargetingSchema>();

  const { control, watch } = methods;
  const individualRulesWatch = watch('individualRules') as IndividualRuleItem[];

  const handleCopyUserId = (value: string) => {
    copyToClipBoard(value);
    notify({
      message: t('message:copied')
    });
  };

  return (
    <Card>
      <div className="flex items-center w-full justify-between">
        <div className="flex items-center gap-x-2">
          <p className="typo-para-medium leading-4 text-gray-700">
            {t('form:targeting.individual-target')}
          </p>
          <Tooltip
            align="start"
            alignOffset={-105}
            content={t('form:targeting.tooltip.individual-target')}
            trigger={
              <div className="flex-center size-fit">
                <Icon icon={IconInfo} size={'xxs'} color="gray-500" />
              </div>
            }
            className="max-w-[400px]"
          />
        </div>
        <div
          className="flex-center h-8 w-8 px-2 border-gray-500 border rounded-md cursor-pointer group"
          onClick={() => handleDiscardChanges(DiscardChangesType.INDIVIDUAL)}
        >
          <Icon
            icon={IconUndoOutlined}
            size={'sm'}
            className="flex-center text-gray-500 group-hover:text-gray-700"
          />
        </div>
      </div>
      {individualRules.map((item, index) => (
        <div key={index} className="flex flex-col w-full gap-y-4">
          <Form.Field
            control={control}
            name={`individualRules.${index}.users`}
            render={({ field }) => {
              return (
                <Form.Item className="py-0">
                  <Form.Label className="flex items-center gap-x-2 !mb-2">
                    <p className="uppercase">{t('feature-flags.serve')}</p>
                    <FlagVariationPolygon index={index} className="!z-0" />
                    <p>{item?.name}</p>
                  </Form.Label>
                  <Form.Control>
                    <div className="flex items-center w-full gap-x-2">
                      <CreatableSelect
                        value={field.value?.map((item: string) => ({
                          label: item,
                          value: item
                        }))}
                        placeholder={t('form:feature-flags.add-user-id')}
                        onChange={(options: MultiValue<Option>) => {
                          const newOption = options.find(o => o['__isNew__']);
                          const alreadyTargetedVariation =
                            getAlreadyTargetedVariation(
                              individualRulesWatch,
                              item.variationId,
                              newOption?.label || ''
                            );
                          if (!alreadyTargetedVariation) {
                            field.onChange(options.map(o => o.value));
                          }
                        }}
                        className="w-full"
                        formatCreateLabel={v => {
                          const isAlreadyExisted = getAlreadyTargetedVariation(
                            individualRulesWatch,
                            item.variationId,
                            v
                          );

                          if (isAlreadyExisted) {
                            const variationName = truncateBySide(
                              isAlreadyExisted.name as string,
                              50
                            );
                            return (
                              <UserMessage
                                message={
                                  <Trans
                                    i18nKey={
                                      'form:feature-flags.value-already-targeted'
                                    }
                                    values={{
                                      value: v,
                                      targetedIn: variationName
                                    }}
                                  />
                                }
                              />
                            );
                          }
                          return (
                            <UserMessage
                              message={t('form:feature-flags.add-user-id')}
                            />
                          );
                        }}
                        noOptionsMessage={({ inputValue }) => {
                          return (
                            <UserMessage
                              message={t(
                                inputValue
                                  ? `form:feature-flags.already-targeted`
                                  : 'form:no-opts-type-to-create'
                              )}
                            />
                          );
                        }}
                      />
                      <Button
                        disabled={!field.value?.length}
                        variant={'secondary-2'}
                        type="button"
                        size={'icon'}
                        tabIndex={-1}
                        onClick={() =>
                          handleCopyUserId(field.value?.join(', '))
                        }
                      >
                        <Icon icon={IconCopy} />
                      </Button>
                    </div>
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              );
            }}
          />
        </div>
      ))}
    </Card>
  );
};

export default IndividualRule;
