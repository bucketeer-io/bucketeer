import React, { LabelHTMLAttributes, Ref } from 'react';
import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import { FormFieldContext, FormItemContext } from 'components/form';

interface FormLabelProps extends LabelHTMLAttributes<HTMLLabelElement> {
  required?: boolean;
  optional?: boolean;
}

const FormLabel = React.forwardRef(
  (
    { className, children, required, optional, ...props }: FormLabelProps,
    ref: Ref<HTMLLabelElement>
  ) => {
    const { t } = useTranslation(['form']);
    const fieldContext = React.useContext(FormFieldContext);
    const itemContext = React.useContext(FormItemContext);
    const formItemId =
      fieldContext && itemContext ? `${itemContext.id}-form-item` : undefined;

    return (
      <label
        ref={ref}
        htmlFor={formItemId}
        className={cn('typo-para-small text-gray-600 mb-1', className)}
        {...props}
      >
        {children}
        {optional && (
          <span className="text-gray-500 ml-2">({t(`optional`)})</span>
        )}
        {required && <span className="text-accent-red-500 ml-0.5">*</span>}
      </label>
    );
  }
);

FormLabel.displayName = 'FormLabel';

export default FormLabel;
