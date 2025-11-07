import { useMemo } from 'react';
import { useFieldArray, useFormContext } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { IconAddOutlined } from 'react-icons-material-design';
import { useTranslation } from 'i18n';
import { v4 as uuid } from 'uuid';
import { cn } from 'utils/style';
import { IconTrash } from '@icons';
import { FlagFormSchema } from 'pages/create-flag/form-schema';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import Button from 'components/button';
import ReactCodeEditor from 'components/code-editor';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';

const Variations = () => {
  const { t } = useTranslation(['form', 'common', 'table']);

  const { control, watch } = useFormContext<FlagFormSchema>();
  const {
    fields: variations,
    append,
    remove
  } = useFieldArray({
    name: 'variations',
    control,
    keyName: 'flagVariation'
  });

  const variationType = watch('variationType');
  const onVariation = watch('defaultOnVariation');
  const offVariation = watch('defaultOffVariation');

  const isBoolean = useMemo(() => variationType === 'BOOLEAN', [variationType]);
  const isJSON = useMemo(() => variationType === 'JSON', [variationType]);
  const isYAML = useMemo(() => variationType === 'YAML', [variationType]);

  const onAddVariation = () => {
    append({
      id: uuid(),
      value: isJSON ? '{}' : '',
      name: '',
      description: ''
    });
  };

  const onDeleteVariation = (itemIndex: number) => {
    remove(itemIndex);
  };

  return (
    <div className="flex flex-col w-full gap-y-4">
      {variations.map((item, variationIndex) => (
        <div key={item.flagVariation} className="flex flex-col w-full">
          <div className="flex items-center gap-x-2 mb-3 typo-para-small text-gray-600">
            <FlagVariationPolygon index={variationIndex} />
            <Trans
              i18nKey={'form:feature-flags.variation'}
              values={{
                index: `${variationIndex + 1}`
              }}
            />
          </div>
          <div className="flex gap-x-4">
            <div className="flex flex-col w-full gap-y-4">
              <div
                className={cn('flex w-full gap-4', {
                  'flex-col': isJSON || isYAML
                })}
              >
                <Form.Field
                  control={control}
                  name={`variations.${variationIndex}.name`}
                  render={({ field }) => (
                    <Form.Item className="w-full py-0">
                      <Form.Label required>{t('common:name')}</Form.Label>
                      <Form.Control>
                        <Input {...field} placeholder={t('placeholder-name')} />
                      </Form.Control>
                      <Form.Message />
                    </Form.Item>
                  )}
                />
                <Form.Field
                  control={control}
                  name={`variations.${variationIndex}.value`}
                  render={({ field }) => (
                    <Form.Item className="w-full py-0">
                      <Form.Label required>
                        {t('feature-flags.value')}
                      </Form.Label>
                      <Form.Control>
                        {isJSON || isYAML ? (
                          <div className="flex flex-col gap-y-2">
                            <ReactCodeEditor
                              defaultLanguage={isYAML ? 'yaml' : 'json'}
                              value={field.value}
                              onChange={field.onChange}
                            />
                          </div>
                        ) : (
                          <Input
                            {...field}
                            placeholder={t(
                              'feature-flags.placeholder-variation'
                            )}
                            disabled={isBoolean}
                          />
                        )}
                      </Form.Control>
                      <Form.Message />
                    </Form.Item>
                  )}
                />
              </div>
              {!isJSON && (
                <Form.Field
                  control={control}
                  name={`variations.${variationIndex}.description`}
                  render={({ field }) => (
                    <Form.Item className="py-0">
                      <Form.Label>{t('description')}</Form.Label>
                      <Form.Control>
                        <Input {...field} placeholder={t('placeholder-desc')} />
                      </Form.Control>
                      <Form.Message />
                    </Form.Item>
                  )}
                />
              )}
            </div>
            <Button
              variant="grey"
              size="icon"
              type="button"
              className="p-0 size-5 mt-5 self-center"
              disabled={
                variations.length <= 2 ||
                [onVariation, offVariation].includes(item.id)
              }
              onClick={() => onDeleteVariation(variationIndex)}
            >
              <Icon icon={IconTrash} size="sm" />
            </Button>
          </div>
        </div>
      ))}
      <Button
        onClick={onAddVariation}
        variant="text"
        size="sm"
        type="button"
        disabled={isBoolean}
        className="w-fit px-0 !typo-para-medium"
      >
        <Icon icon={IconAddOutlined} size="sm" />
        {t(`feature-flags.add-variation`)}
      </Button>
    </div>
  );
};

export default Variations;
