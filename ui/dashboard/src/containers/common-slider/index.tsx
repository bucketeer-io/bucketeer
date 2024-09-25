import { ReactNode } from 'react';
import CommonForm, { CommonFormProps } from 'containers/common-form';
import { useTranslation } from 'i18n';
import { Button } from 'components/button';
import { ButtonBar } from 'components/button-bar';
import SlideModal, { SliderProps } from 'components/modal/slide';

type CommonSliderProps = CommonFormProps &
  SliderProps & {
    submitTextBtn: string;
    children?: ReactNode;
    isLoading?: boolean;
  };

const CommonSlider = ({
  isOpen,
  formFields,
  submitTextBtn,
  children,
  title,
  isLoading,
  formSchema,
  formData,
  onClose,
  onSubmit
}: CommonSliderProps) => {
  const { t } = useTranslation(['form']);

  const handleOnClickSubmitBtn = () => {
    document.getElementById('common-submit-btn')?.click();
  };

  return (
    <SlideModal title={title} isOpen={isOpen} onClose={onClose}>
      <div className="p-5 flex flex-col gap-5">
        <h3 className="typo-head-bold-small text-gray-800">
          {t(`general-info`)}
        </h3>
        <CommonForm
          className="p-0 shadow-none"
          formFields={formFields}
          isShowSubmitBtn={false}
          formClassName="mt-0"
          formSchema={formSchema}
          formData={formData}
          onSubmit={onSubmit}
        />
        <>{children}</>
      </div>

      <div className="absolute bottom-0 bg-gray-50 w-full rounded-b-lg">
        <ButtonBar
          primaryButton={
            <Button variant="secondary" onClick={onClose}>
              {t(`common:cancel`)}
            </Button>
          }
          secondaryButton={
            <Button loading={isLoading} onClick={handleOnClickSubmitBtn}>
              {submitTextBtn}
            </Button>
          }
        />
      </div>
    </SlideModal>
  );
};

export default CommonSlider;
