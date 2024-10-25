import React, { HTMLAttributes, Ref } from 'react';
import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import { useFormField } from 'components/form';

interface FormLabelProps extends HTMLAttributes<HTMLDivElement> {
  required?: boolean;
  optional?: boolean;
}

const FormLabel = React.forwardRef(
  (
    { className, children, required, optional, ...props }: FormLabelProps,
    ref: Ref<HTMLDivElement>
  ) => {
    const { t } = useTranslation(['form']);
    const { formItemId } = useFormField();

    return (
      <div
        ref={ref}
        className={cn('typo-para-small text-gray-600 mb-1', className)}
        id={formItemId}
        {...props}
      >
        {children}
        {optional && (
          <span className="text-gray-500 ml-2">({t(`optional`)})</span>
        )}
        {required && <span className="text-accent-red-500 ml-0.5">*</span>}
      </div>
    );
  }
);

export default FormLabel;
