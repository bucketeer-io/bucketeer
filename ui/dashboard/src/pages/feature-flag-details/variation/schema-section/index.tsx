import { useMemo } from 'react';
import { useFormContext } from 'react-hook-form';
import { useToggleOpen } from 'hooks';
import { useTranslation } from 'i18n';
import { Feature } from '@types';
import { isSchemaSupported } from 'utils/variation-value-schema';
import Button from 'components/button';
import ReactCodeEditor from 'components/code-editor';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import Card from '../../elements/card';
import { VariationForm } from '../form-schema';
import SchemaDialog from './schema-dialog';

const SchemaSection = ({
  feature,
  editable
}: {
  feature: Feature;
  editable: boolean;
}) => {
  const { t } = useTranslation(['common', 'form']);
  const { watch, setValue, trigger } = useFormContext<VariationForm>();
  const [isDialogOpen, onOpenDialog, onCloseDialog] = useToggleOpen(false);

  const schema = watch('variationValueSchema');

  const typeLabel = useMemo(() => {
    switch (schema?.type) {
      case 'ENUM':
        return t('form:feature-flags.value-schema.type-enum');
      case 'REGEX':
        return t('form:feature-flags.value-schema.type-regex');
      case 'JSON_SCHEMA':
        return t('form:feature-flags.value-schema.type-json-schema');
      default:
        return '';
    }
  }, [schema, t]);

  const formattedJsonSchema = useMemo(() => {
    const value = schema?.jsonSchemaValidator?.schema;
    if (!value) return '';
    try {
      return JSON.stringify(JSON.parse(value), null, 2);
    } catch {
      return value;
    }
  }, [schema]);

  if (!isSchemaSupported(feature.variationType)) return null;

  const handleRemove = () => {
    setValue('variationValueSchema', null, {
      shouldDirty: true,
      shouldValidate: true
    });
    trigger('variations');
  };

  return (
    <Card className="gap-y-4">
      <div className="flex items-center justify-between w-full gap-x-6">
        <div className="flex flex-col gap-y-1">
          <div className="flex items-center gap-x-2">
            <h3 className="typo-head-bold-small text-gray-800">
              {t('form:feature-flags.value-schema.title')}
            </h3>
            {schema && (
              <span className="px-2 py-0.5 rounded bg-primary-50 text-primary-500 typo-para-small">
                {typeLabel}
              </span>
            )}
          </div>
          <p className="typo-para-small text-gray-600">
            {schema
              ? schema.description ||
                t('form:feature-flags.value-schema.has-schema')
              : t('form:feature-flags.value-schema.no-schema')}
          </p>
        </div>
        <div className="flex items-center gap-x-3 shrink-0">
          {schema ? (
            <>
              <DisabledButtonTooltip
                hidden={editable}
                trigger={
                  <Button
                    type="button"
                    variant="grey"
                    disabled={!editable}
                    onClick={handleRemove}
                  >
                    {t('form:feature-flags.value-schema.remove')}
                  </Button>
                }
              />
              <DisabledButtonTooltip
                hidden={editable}
                trigger={
                  <Button
                    type="button"
                    variant="secondary"
                    disabled={!editable}
                    onClick={onOpenDialog}
                  >
                    {t('form:feature-flags.value-schema.edit')}
                  </Button>
                }
              />
            </>
          ) : (
            <DisabledButtonTooltip
              hidden={editable}
              trigger={
                <Button
                  type="button"
                  variant="secondary"
                  disabled={!editable}
                  onClick={onOpenDialog}
                >
                  {t('form:feature-flags.value-schema.add')}
                </Button>
              }
            />
          )}
        </div>
      </div>

      {schema?.type === 'ENUM' && (
        <div className="flex flex-wrap gap-2">
          {schema.enumValidator?.values?.map(value => (
            <span
              key={value}
              className="px-2 py-1 rounded bg-gray-100 text-gray-700 typo-para-small"
            >
              {value}
            </span>
          ))}
        </div>
      )}
      {schema?.type === 'REGEX' && (
        <code className="w-fit px-3 py-2 rounded-lg bg-gray-100 text-gray-700 typo-para-small">
          {schema.regexValidator?.pattern}
        </code>
      )}
      {schema?.type === 'JSON_SCHEMA' && (
        <ReactCodeEditor
          defaultLanguage="json"
          value={formattedJsonSchema}
          readOnly
          isExpand={false}
          className="min-h-[120px] h-[120px]"
        />
      )}

      {isDialogOpen && (
        <SchemaDialog
          isOpen={isDialogOpen}
          feature={feature}
          onClose={onCloseDialog}
        />
      )}
    </Card>
  );
};

export default SchemaSection;
