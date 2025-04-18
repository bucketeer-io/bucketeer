import { memo } from 'react';
import { useTranslation } from 'i18n';
import Button from 'components/button';

interface Props {
  onApply: () => void;
  onCancel: () => void;
}

const ActionBar = memo(({ onApply, onCancel }: Props) => {
  const { t } = useTranslation(['common', 'form']);

  return (
    <div className="sticky bottom-0 left-0 right-0 flex items-center justify-end w-full gap-x-4 p-5 border-t border-gray-200 bg-white">
      <Button variant="secondary" onClick={onCancel}>
        {t('cancel')}
      </Button>
      <Button onClick={onApply}>{t('apply')}</Button>
    </div>
  );
});

export default ActionBar;
