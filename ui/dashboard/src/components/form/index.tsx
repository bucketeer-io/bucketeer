import * as React from 'react';
import { FieldPath, FieldValues, useFormContext } from 'react-hook-form';
import { cn } from 'utils/style';
import FormControl from './form-control';
import FormField from './form-field';
import FormItem from './form-item';
import FormLabel from './form-label';
import FormMessage from './form-message';

type FormItemContextValue = {
  id: string;
};

export const FormItemContext = React.createContext<FormItemContextValue | null>(
  null
);

type FormFieldContextValue<
  TFieldValues extends FieldValues = FieldValues,
  TName extends FieldPath<TFieldValues> = FieldPath<TFieldValues>
> = {
  name: TName;
};

export const FormFieldContext =
  React.createContext<FormFieldContextValue | null>(null);

export const useFormField = () => {
  const fieldContext = React.useContext(FormFieldContext);
  const itemContext = React.useContext(FormItemContext);

  if (!fieldContext) {
    throw new Error('useFormField should be used within <FormField>');
  }

  if (!itemContext) {
    throw new Error('useFormField should be used within <FormItem>');
  }

  const { getFieldState, formState } = useFormContext();
  const fieldState = getFieldState(fieldContext.name, formState);
  const { id } = itemContext;

  return {
    id,
    name: fieldContext.name,
    formItemId: `${id}-form-item`,
    formDescriptionId: `${id}-form-item-description`,
    formMessageId: `${id}-form-item-message`,
    ...fieldState
  };
};

interface FormProps extends React.FormHTMLAttributes<HTMLFormElement> {
  children: React.ReactNode;
  onSubmit: (e: React.FormEvent<HTMLFormElement>) => void;
  className?: string;
}

const Form = ({
  children,
  onSubmit,
  className,
  autoComplete = 'on',
  ...props
}: FormProps) => {
  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    onSubmit(e);
  };

  return (
    <form
      onSubmit={handleSubmit}
      className={cn(className)}
      autoComplete={autoComplete}
      {...props}
    >
      {children}
    </form>
  );
};

Form.displayName = 'Form';
Form.Control = FormControl;
Form.Label = FormLabel;
Form.Item = FormItem;
Form.Message = FormMessage;
Form.Field = FormField;

export default Form;
