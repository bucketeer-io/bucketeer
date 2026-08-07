import { useCallback, useMemo, useState } from 'react';
import { useFormContext } from 'react-hook-form';
import { useTranslation } from 'i18n';
import {
  Feature,
  VariationValueSchema,
  VariationValueSchemaType
} from '@types';
import {
  createValueValidator,
  getSupportedSchemaTypes,
  validateSchemaDefinition
} from 'utils/variation-value-schema';
import { IconAlert } from '@icons';
import Button from 'components/button';
import ReactCodeEditor from 'components/code-editor';
import { CreatableSelect } from 'components/creatable-select';
import Icon from 'components/icon';
import Input from 'components/input';
import DialogModal from 'components/modal/dialog';
import { RadioGroup, RadioGroupItem } from 'components/radio';
import { VariationForm } from '../form-schema';

interface SchemaDialogProps {
  isOpen: boolean;
  feature: Feature;
  onClose: () => void;
}

interface TestResult {
  index: number;
  label: string;
  passed: boolean;
}

const DEFAULT_JSON_SCHEMA = `{
  "type": "object",
  "properties": {}
}`;

const SchemaDialog = ({ isOpen, feature, onClose }: SchemaDialogProps) => {
  const { t } = useTranslation(['common', 'form', 'message']);
  const { watch, setValue, trigger } = useFormContext<VariationForm>();

  const currentSchema = watch('variationValueSchema');
  const variations = watch('variations');

  const supportedTypes = useMemo(
    () => getSupportedSchemaTypes(feature.variationType),
    [feature.variationType]
  );

  const [schemaType, setSchemaType] = useState<VariationValueSchemaType>(
    currentSchema?.type && supportedTypes.includes(currentSchema.type)
      ? currentSchema.type
      : supportedTypes[0]
  );
  const [description, setDescription] = useState(
    currentSchema?.description ?? ''
  );
  const [enumValues, setEnumValues] = useState<string[]>(
    currentSchema?.enumValidator?.values ?? []
  );
  const [regexPattern, setRegexPattern] = useState(
    currentSchema?.regexValidator?.pattern ?? ''
  );
  const [jsonSchema, setJsonSchema] = useState(
    currentSchema?.jsonSchemaValidator?.schema ?? DEFAULT_JSON_SCHEMA
  );
  const [testResults, setTestResults] = useState<TestResult[] | null>(null);

  const buildSchema = useCallback((): VariationValueSchema => {
    const schema: VariationValueSchema = { type: schemaType };
    if (description.trim()) schema.description = description.trim();
    if (schemaType === 'ENUM') schema.enumValidator = { values: enumValues };
    if (schemaType === 'REGEX')
      schema.regexValidator = { pattern: regexPattern };
    if (schemaType === 'JSON_SCHEMA')
      schema.jsonSchemaValidator = { schema: jsonSchema };
    return schema;
  }, [schemaType, description, enumValues, regexPattern, jsonSchema]);

  const definitionError = useMemo(
    () => validateSchemaDefinition(buildSchema(), feature.variationType),
    [buildSchema, feature.variationType]
  );

  const definitionErrorMessage = useMemo(() => {
    switch (definitionError) {
      case 'enum-not-number':
        return t('message:validation.value-schema-enum-not-number');
      case 'regex-invalid':
        return t('message:validation.value-schema-invalid-regex');
      case 'json-schema-invalid':
        return t('message:validation.value-schema-invalid-json-schema');
      // Empty-field errors only disable the buttons; no message needed.
      default:
        return null;
    }
  }, [definitionError, t]);

  const handleTest = useCallback(() => {
    const validator = createValueValidator(
      buildSchema(),
      feature.variationType
    );
    if (!validator) return;
    setTestResults(
      variations.map((variation, index) => ({
        index,
        label:
          variation.name ||
          variation.value ||
          t('form:feature-flags.variation', { index: index + 1 }),
        passed: validator(variation.value)
      }))
    );
  }, [buildSchema, feature.variationType, variations, t]);

  const handleSave = useCallback(() => {
    setValue('variationValueSchema', buildSchema(), {
      shouldDirty: true,
      shouldValidate: true
    });
    trigger('variations');
    onClose();
  }, [buildSchema, setValue, trigger, onClose]);

  const failedCount = testResults?.filter(item => !item.passed).length ?? 0;

  return (
    <DialogModal
      className="w-[600px]"
      title={
        currentSchema
          ? t('form:feature-flags.value-schema.edit-title')
          : t('form:feature-flags.value-schema.add-title')
      }
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="flex flex-col w-full gap-y-5 p-5">
        {supportedTypes.length > 1 && (
          <div className="flex flex-col gap-y-3">
            <p className="typo-para-small text-gray-600">
              {t('form:feature-flags.value-schema.type')}
            </p>
            <RadioGroup
              value={schemaType}
              onValueChange={value => {
                setSchemaType(value as VariationValueSchemaType);
                setTestResults(null);
              }}
              className="flex items-center gap-x-6"
            >
              {supportedTypes.map(type => (
                <label
                  key={type}
                  className="flex items-center gap-x-2 cursor-pointer typo-para-medium text-gray-700"
                >
                  <RadioGroupItem value={type} id={`schema-type-${type}`} />
                  {t(
                    `form:feature-flags.value-schema.type-${
                      type === 'JSON_SCHEMA'
                        ? 'json-schema'
                        : type.toLowerCase()
                    }`
                  )}
                </label>
              ))}
            </RadioGroup>
          </div>
        )}

        {schemaType === 'ENUM' && (
          <div className="flex flex-col gap-y-2">
            <p className="typo-para-small text-gray-600">
              {t('form:feature-flags.value-schema.enum-values')}
            </p>
            <CreatableSelect
              value={enumValues.map(value => ({ label: value, value }))}
              placeholder={t(
                'form:feature-flags.value-schema.enum-values-placeholder'
              )}
              onChange={options => {
                setEnumValues(options.map(option => option.value));
                setTestResults(null);
              }}
              noOptionsMessage={() => null}
            />
            {feature.variationType === 'NUMBER' && (
              <p className="typo-para-small text-gray-500">
                {t('form:feature-flags.value-schema.enum-values-number-note')}
              </p>
            )}
          </div>
        )}

        {schemaType === 'REGEX' && (
          <div className="flex flex-col gap-y-2">
            <p className="typo-para-small text-gray-600">
              {t('form:feature-flags.value-schema.regex-pattern')}
            </p>
            <Input
              value={regexPattern}
              placeholder={t(
                'form:feature-flags.value-schema.regex-pattern-placeholder'
              )}
              onChange={value => {
                setRegexPattern(value);
                setTestResults(null);
              }}
            />
            <p className="typo-para-small text-gray-500">
              {t('form:feature-flags.value-schema.regex-note')}
            </p>
          </div>
        )}

        {schemaType === 'JSON_SCHEMA' && (
          <div className="flex flex-col gap-y-2">
            <p className="typo-para-small text-gray-600">
              {t('form:feature-flags.value-schema.json-schema')}
            </p>
            <ReactCodeEditor
              defaultLanguage="json"
              value={jsonSchema}
              onChange={value => {
                setJsonSchema(value ?? '');
                setTestResults(null);
              }}
              isExpand={false}
              className="min-h-[240px] h-[240px]"
            />
            <p className="typo-para-small text-gray-500">
              {t('form:feature-flags.value-schema.json-schema-note')}
            </p>
          </div>
        )}

        <div className="flex flex-col gap-y-2">
          <p className="typo-para-small text-gray-600">
            {t('form:description')}
          </p>
          <Input
            value={description}
            placeholder={t(
              'form:feature-flags.value-schema.description-placeholder'
            )}
            onChange={value => setDescription(value)}
          />
        </div>

        {definitionErrorMessage && (
          <div className="flex items-center gap-x-2 typo-para-small text-accent-red-500">
            <Icon icon={IconAlert} size="xs" />
            {definitionErrorMessage}
          </div>
        )}

        {testResults && (
          <div className="flex flex-col gap-y-2 p-3 rounded-lg bg-gray-100">
            <p
              className={
                failedCount > 0
                  ? 'typo-para-small text-accent-red-500'
                  : 'typo-para-small text-accent-green-500'
              }
            >
              {failedCount > 0
                ? t('form:feature-flags.value-schema.test-fail', {
                    count: failedCount
                  })
                : t('form:feature-flags.value-schema.test-pass')}
            </p>
            <ul className="flex flex-col gap-y-1">
              {testResults.map(result => (
                <li
                  key={result.index}
                  className="flex items-center justify-between typo-para-small text-gray-700"
                >
                  <span className="truncate max-w-[400px]">{result.label}</span>
                  <span
                    className={
                      result.passed
                        ? 'text-accent-green-500'
                        : 'text-accent-red-500'
                    }
                  >
                    {result.passed
                      ? t('form:feature-flags.value-schema.test-valid')
                      : t('form:feature-flags.value-schema.test-invalid')}
                  </span>
                </li>
              ))}
            </ul>
          </div>
        )}

        <div className="flex items-center justify-end gap-x-3 pt-2">
          <Button type="button" variant="secondary" onClick={onClose}>
            {t('common:cancel')}
          </Button>
          <Button
            type="button"
            variant="secondary"
            disabled={!!definitionError}
            onClick={handleTest}
          >
            {t('form:feature-flags.value-schema.test')}
          </Button>
          <Button
            type="button"
            disabled={!!definitionError}
            onClick={handleSave}
          >
            {t('common:save')}
          </Button>
        </div>
      </div>
    </DialogModal>
  );
};

export default SchemaDialog;
