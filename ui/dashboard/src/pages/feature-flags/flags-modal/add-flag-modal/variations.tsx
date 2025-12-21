import { Fragment, useMemo, useState } from 'react';
import { useFieldArray, useFormContext } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { IconAddOutlined } from 'react-icons-material-design';
import { useTranslation } from 'i18n';
import { v4 as uuid } from 'uuid';
import { FeatureVariationType } from '@types';
import { cn } from 'utils/style';
import { IconTrash } from '@icons';
import { FlagFormSchema } from 'pages/create-flag/form-schema';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import Button from 'components/button';
import ReactCodeEditor from 'components/code-editor';
import ReactCodeEditorModal from 'components/code-editor/code-editor-mode';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';

const Variations = ({
  variationType,
  refModel
}: {
  variationType: FeatureVariationType;
  refModel?: React.RefObject<HTMLDivElement>;
}) => {
  const { t } = useTranslation(['common', 'form']);

  const methods = useFormContext<FlagFormSchema>();
  const {
    control,
    watch,
    trigger,
    formState: { errors }
  } = methods;

  const [variationSelected, setVariationSelected] = useState<number | null>(
    null
  );
  const {
    fields: variations,
    append,
    remove
  } = useFieldArray({
    name: 'variations',
    control,
    keyName: 'flagVariation'
  });

  const isBoolean = useMemo(() => variationType === 'BOOLEAN', [variationType]);
  const isJSON = useMemo(() => variationType === 'JSON', [variationType]);
  const isYAML = useMemo(() => variationType === 'YAML', [variationType]);
  const onVariation = watch('defaultOnVariation');
  const offVariation = watch('defaultOffVariation');

  const onAddVariation = () => {
    append({
      id: uuid(),
      value: '',
      name: '',
      description: ''
    });
  };

  const onDeleteVariation = (itemIndex: number) => {
    remove(itemIndex);
  };

  const revalidateVariations = () => {
    if (!errors || !errors?.variations || !errors.variations.length) return;
    errors.variations.map?.((_, index) => {
      trigger(`variations.${index}.value`);
    });
  };
  return (
    <>
      {variations.map((item, variationIndex) => (
        <Fragment key={item.flagVariation}>
          <div className="flex items-center gap-x-2 mb-3 typo-para-small text-gray-600">
            <FlagVariationPolygon index={variationIndex} />
            <Trans
              i18nKey={'form:feature-flags.variation'}
              values={{
                index: `${variationIndex + 1}`
              }}
            />
          </div>
          <div
            key={item.flagVariation}
            className={cn('flex w-full gap-x-4 mb-3 ')}
          >
            <div
              className={cn('w-full flex gap-x-4 [&>div]:flex-1', {
                'flex-col': isJSON || isYAML
              })}
            >
              <Form.Field
                control={control}
                name={`variations.${variationIndex}.name`}
                render={({ field }) => (
                  <Form.Item className="py-0">
                    <Form.Label required>{t('name')}</Form.Label>
                    <Form.Control>
                      <Input
                        {...field}
                        placeholder={t('form:placeholder-name')}
                        name="flag-variation-name"
                      />
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />

              <Form.Field
                control={control}
                name={`variations.${variationIndex}.value`}
                render={({ field }) => (
                  <Form.Item
                    className={cn('py-0', { 'py-3': isJSON || isYAML })}
                  >
                    <div className="flex items-center gap-x-2">
                      <Form.Label required>
                        {t('form:feature-flags.value')}
                      </Form.Label>
                    </div>
                    <Form.Control>
                      {isJSON || isYAML ? (
                        <div className="flex flex-col gap-y-2">
                          <div className="w-full max-h-[100px] overflow-hidden">
                            <ReactCodeEditor
                              scrollParent={refModel}
                              defaultLanguage={isYAML ? 'yaml' : 'json'}
                              value={field.value}
                              onChange={value => {
                                field.onChange(value);
                                revalidateVariations();
                              }}
                              onExpand={() =>
                                setVariationSelected(variationIndex)
                              }
                              className="!min-h-[100px] h-full"
                              lastLine={2}
                            />
                          </div>
                          <ReactCodeEditorModal
                            defaultLanguage={isYAML ? 'yaml' : 'json'}
                            isOpen={variationSelected === variationIndex}
                            onClose={() => {
                              setVariationSelected(null);
                            }}
                            title={
                              <div className="flex items-center gap-x-2 typo-para-big text-gray-600 font-bold">
                                <FlagVariationPolygon index={variationIndex} />
                                <Trans
                                  i18nKey={'form:feature-flags.variation'}
                                  values={{
                                    index: `${variationIndex + 1}`
                                  }}
                                />
                              </div>
                            }
                            value={field.value}
                            onChange={value => {
                              field.onChange(value);
                              revalidateVariations();
                            }}
                            error={
                              errors.variations?.[variationIndex]?.value
                                ?.message as string
                            }
                          />
                        </div>
                      ) : (
                        <Input
                          {...field}
                          placeholder={t(
                            'form:feature-flags.placeholder-variation'
                          )}
                          onChange={value => {
                            field.onChange(value);
                            revalidateVariations();
                          }}
                          disabled={isBoolean}
                          className={isBoolean ? 'capitalize' : ''}
                        />
                      )}
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />
            </div>
            {variations.length > 2 &&
              ![onVariation, offVariation].includes(item.id) && (
                <Button
                  variant="text"
                  size="icon"
                  type="button"
                  className="p-0 size-5 mt-5 self-center"
                  onClick={() => onDeleteVariation(variationIndex)}
                >
                  <Icon icon={IconTrash} size="sm" color="gray-600" />
                </Button>
              )}
          </div>
        </Fragment>
      ))}
      <Button
        onClick={onAddVariation}
        variant="text"
        type="button"
        disabled={isBoolean}
        className="my-1"
      >
        <Icon icon={IconAddOutlined} />
        {t(`form:feature-flags.add-variation`)}
      </Button>
    </>
  );
};

export default Variations;
