import { ReactNode, useCallback, useMemo } from 'react';
import { useFormContext } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { MultiValue } from 'react-select';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import { cloneDeep } from 'lodash';
import { truncateBySide } from 'utils/converts';
import { copyToClipBoard } from 'utils/function';
import { IconCopy, IconInfo } from '@icons';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import Button from 'components/button';
import { CreatableSelect, Option } from 'components/creatable-select';
import Form from 'components/form';
import Icon from 'components/icon';
import Card from '../../elements/card';
import { IndividualRuleItem } from '../types';
import { getAlreadyTargetedVariation } from '../utils';

interface Props {
  individualRules: IndividualRuleItem[];
  onChangeIndividualRules: (value: IndividualRuleItem[]) => void;
}

const UserMessage = ({ message }: { message: ReactNode }) => {
  return <div className={'text-center text-gray-500'}>{message}</div>;
};

const IndividualRule = ({
  individualRules,
  onChangeIndividualRules
}: Props) => {
  const { t } = useTranslation(['table', 'form', 'common']);
  const { notify } = useToast();

  const cloneIndividualRules = useMemo(
    () => cloneDeep(individualRules),
    [individualRules]
  );
  const methods = useFormContext();

  const { control } = methods;

  const onChangeFormField = useCallback(
    (individualIndex: number, field: string, value: string[]) => {
      cloneIndividualRules[individualIndex] = {
        ...cloneIndividualRules[individualIndex],
        [field]: value
      };
      return onChangeIndividualRules(cloneIndividualRules);
    },
    [cloneIndividualRules]
  );

  const handleCopyUserId = (value: string) => {
    copyToClipBoard(value);
    notify({
      toastType: 'toast',
      messageType: 'success',
      message: 'Copied'
    });
  };

  return (
    individualRules.length > 0 && (
      <Card>
        <div>
          <div className="flex items-center gap-x-2">
            <p className="typo-para-medium leading-4 text-gray-700">
              {t('form:feature-flags.prerequisites')}
            </p>
            <Icon icon={IconInfo} size={'xxs'} color="gray-500" />
          </div>
        </div>
        {individualRules.map((item, index) => (
          <div key={index} className="flex flex-col w-full gap-y-4">
            <Form.Field
              control={control}
              name={`targetIndividualRules.${index}.users`}
              render={({ field }) => {
                return (
                  <Form.Item className="py-0">
                    <Form.Label
                      required
                      className="flex items-center gap-x-2 !mb-2"
                    >
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
                                individualRules,
                                item.variationId,
                                newOption?.label || ''
                              );
                            if (!alreadyTargetedVariation) {
                              field.onChange(options.map(o => o.value));
                              onChangeFormField(
                                index,
                                'users',
                                options.map(o => o.value)
                              );
                            }
                          }}
                          className="w-full"
                          formatCreateLabel={v => {
                            const isAlreadyExisted =
                              getAlreadyTargetedVariation(
                                individualRules,
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
                                      i18nKey={'form:feature-flags.add-user-id'}
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
                                  `form:feature-flags.${inputValue ? 'already-targeted' : 'add-user-id'}`
                                )}
                              />
                            );
                          }}
                        />
                        <Button
                          disabled={!field.value?.length}
                          variant={'secondary-2'}
                          size={'icon'}
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
    )
  );
};

export default IndividualRule;
