import { useMemo } from 'react';
import { useFormContext } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { IconAddOutlined } from 'react-icons-material-design';
import { useTranslation } from 'i18n';
import { v4 as uuid } from 'uuid';
import { FeatureVariation, FeatureVariationType } from '@types';
import { IconTrash } from '@icons';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import Button from 'components/button';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import TextArea from 'components/textarea';

const Variations = ({
  variationType,
  variations,
  onChangeVariations
}: {
  variationType: FeatureVariationType;
  variations: FeatureVariation[];
  onChangeVariations: (v: FeatureVariation[]) => void;
}) => {
  const { t } = useTranslation(['common', 'form']);

  const methods = useFormContext();
  const { control, watch } = methods;

  const isBoolean = useMemo(() => variationType === 'BOOLEAN', [variationType]);
  const isJSON = useMemo(() => variationType === 'JSON', [variationType]);

  const onVariation = watch('defaultOnVariation');
  const offVariation = watch('defaultOffVariation');

  const onAddVariation = () => {
    onChangeVariations([
      ...variations,
      {
        id: uuid(),
        value: '',
        name: '',
        description: ''
      }
    ]);
  };

  const onDeleteVariation = (itemIndex: number) => {
    const _variations = variations.filter(
      (_item, index) => itemIndex !== index
    );
    onChangeVariations([..._variations]);
  };

  return (
    <>
      {variations.map((item, variationIndex) => (
        <div key={variationIndex} className="flex flex-col w-full">
          <Form.Field
            control={control}
            name={`variations.${variationIndex}.value`}
            render={({ field }) => (
              <Form.Item className="py-2">
                <div className="flex items-center gap-x-2 mb-1">
                  <FlagVariationPolygon index={variationIndex} />
                  <Form.Label required>
                    <Trans
                      i18nKey={'form:feature-flags.variation'}
                      values={{
                        index: `${variationIndex + 1}`
                      }}
                    />
                  </Form.Label>
                </div>
                <Form.Control>
                  {isJSON ? (
                    <TextArea
                      {...field}
                      placeholder={t(
                        'form:feature-flags.placeholder-variation'
                      )}
                      value={item.value}
                    />
                  ) : (
                    <Input
                      {...field}
                      placeholder={t(
                        'form:feature-flags.placeholder-variation'
                      )}
                      disabled={isBoolean}
                      value={item.value}
                      className={isBoolean ? 'capitalize' : ''}
                    />
                  )}
                </Form.Control>
                <Form.Message />
              </Form.Item>
            )}
          />
          <div className="flex items-end mt-5 gap-x-4">
            <div className="flex flex-col flex-1 gap-y-5 pl-4 border-l border-primary-500">
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
                        value={item.name}
                      />
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />
              <Form.Field
                control={control}
                name={`variations.${variationIndex}.description`}
                render={({ field }) => (
                  <Form.Item className="py-0">
                    <Form.Label>{t('form:description')}</Form.Label>
                    <Form.Control>
                      <TextArea
                        {...field}
                        placeholder={t('form:placeholder-desc')}
                        value={item.description}
                      />
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
        </div>
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
