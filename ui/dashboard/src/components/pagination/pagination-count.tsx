import { useTranslation } from 'i18n';

export type PaginationCountProps = {
  totalItems: number;
  value?: number;
};

const PaginationCount = ({ totalItems, value = 0 }: PaginationCountProps) => {
  const { t } = useTranslation(['common']);
  return (
    <p className="text-gray-600 typo-para-medium">
      {t('pagination-count', {
        count: value,
        total: totalItems
      })}
    </p>
  );
};

export default PaginationCount;
