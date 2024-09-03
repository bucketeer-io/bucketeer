import React, { HTMLAttributes, Ref } from 'react';
import { cn } from 'utils/style';
import { useFormField } from 'components/form';

interface FormLabelProps extends HTMLAttributes<HTMLDivElement> {
  required?: boolean;
}

const FormLabel = React.forwardRef(
  (
    { className, children, required, ...props }: FormLabelProps,
    ref: Ref<HTMLDivElement>
  ) => {
    const { formItemId } = useFormField();

    return (
      <div
        ref={ref}
        className={cn('typo-para-small text-gray-600 mb-1', className)}
        id={formItemId}
        {...props}
      >
        {children}
        {required && <span className="text-accent-red-500 ml-0.5">*</span>}
      </div>
    );
  }
);

export default FormLabel;
