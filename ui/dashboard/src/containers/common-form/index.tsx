import { ReactNode, useMemo } from 'react';
import {
  ControllerRenderProps,
  FieldValues,
  FormProvider,
  useForm
} from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { useTranslation } from 'i18n';
import { AnyObject, ObjectSchema } from 'yup';
import { cn } from 'utils/style';
import { Button } from 'components/button';
import Card from 'components/card';
import {
  DropdownMenuContent,
  DropdownMenu,
  DropdownOption,
  DropdownMenuTrigger,
  DropdownMenuItem
} from 'components/dropdown';
import Form from 'components/form';
import Input from 'components/input';
import TextArea from 'components/textarea';

type FieldType = 'input' | 'textarea' | 'dropdown' | 'additional';

export type FormFieldProps = {
  name: string;
  label?: string;
  labelIcon?: ReactNode;
  labelTooltipContent?: string;
  placeholder?: string;
  isRequired?: boolean;
  isOptional?: boolean;
  fieldType?: FieldType;
  isExpand?: boolean;
  isReadOnly?: boolean;
  isDisabled?: boolean;
  dropdownOptions?: DropdownOption[];
  render?: (field: ControllerRenderProps<FieldValues, string>) => ReactNode;
  renderTrigger?: (value: string | number | boolean) => ReactNode;
};

export type CommonFormProps = {
  title?: string;
  formFields: FormFieldProps[];
  formClassName?: string;
  isShowSubmitBtn?: boolean;
  className?: string;
  formSchema?: ObjectSchema<AnyObject>;
  formData?: FieldValues;
  isLoading?: boolean;
  onSubmit: (values: FieldValues) => void;
};

export type FormValues = {
  [key: string]: string | boolean;
};

const CommonForm = ({
  title,
  formFields,
  formClassName,
  isShowSubmitBtn = true,
  className,
  formSchema = {} as ObjectSchema<AnyObject>,
  formData,
  isLoading,
  onSubmit
}: CommonFormProps) => {
  const { t } = useTranslation(['common']);

  const defaultValues: FieldValues = useMemo(() => {
    const obj = {};
    formFields.forEach(item =>
      Object.assign(obj, {
        [item.name]: ''
      })
    );
    return obj;
  }, [formFields]);

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues,
    values: formData
  });

  const renderField = (
    item: FormFieldProps,
    field: ControllerRenderProps<FieldValues, string>,
    fieldType?: FieldType
  ) => {
    switch (fieldType) {
      case 'textarea':
        return <TextArea placeholder={item.placeholder} rows={4} {...field} />;
      case 'dropdown':
        return (
          <DropdownMenu>
            <DropdownMenuTrigger
              disabled={item?.isDisabled}
              label={
                field?.value
                  ? item?.dropdownOptions?.find(
                      opt => opt?.value === field?.value
                    )?.label
                  : item.label
              }
              variant="secondary"
              isExpand={item.isExpand}
              trigger={item?.renderTrigger && item.renderTrigger(field?.value)}
            />
            <div className="relative w-full">
              <DropdownMenuContent
                align="start"
                isExpand={item.isExpand}
                {...field}
              >
                {item?.dropdownOptions && item?.dropdownOptions?.length > 0 ? (
                  item?.dropdownOptions?.map((item, index) => (
                    <DropdownMenuItem
                      {...field}
                      key={index}
                      value={item.value}
                      label={item.label}
                      onSelectOption={value => {
                        field.onChange(value);
                      }}
                    />
                  ))
                ) : (
                  <DropdownMenuItem
                    label="No Result"
                    value="no result"
                    disabled={true}
                  />
                )}
              </DropdownMenuContent>
            </div>
          </DropdownMenu>
        );
      case 'additional':
        return item?.render ? item.render(field) : <></>;
      default:
        return (
          <Input
            placeholder={item.placeholder}
            readOnly={item?.isReadOnly}
            disabled={item?.isDisabled}
            {...field}
          />
        );
    }
  };

  return (
    <Card className={cn('p-5', className)}>
      {title && (
        <p className="text-gray-800 typo-head-bold-small">{`${t(title)}`}</p>
      )}
      <FormProvider {...form}>
        <Form
          onSubmit={form.handleSubmit(onSubmit)}
          className={cn('grid grid-cols-2 gap-5 mt-5', formClassName)}
        >
          {formFields.map((item, index) => (
            <Form.Field
              key={index}
              control={form.control}
              name={item.name}
              render={({ field }) => (
                <Form.Item
                  className={cn('py-0', {
                    'col-span-2': item.isExpand
                  })}
                >
                  <Form.Label
                    optional={item?.isOptional}
                    required={item?.isRequired}
                    tooltipIcon={item?.labelIcon}
                    tooltipContent={item?.labelTooltipContent}
                  >
                    {item.label}
                  </Form.Label>
                  <Form.Control>
                    {renderField(item, field, item?.fieldType)}
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
          ))}
          <Button
            loading={isLoading}
            id="common-submit-btn"
            type="submit"
            className={cn('w-fit', {
              hidden: !isShowSubmitBtn
            })}
          >
            Save
          </Button>
        </Form>
      </FormProvider>
    </Card>
  );
};

export default CommonForm;
