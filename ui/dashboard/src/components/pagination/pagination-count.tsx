import { useTranslation } from 'i18n';

export type PaginationCountProps = {
  totalItems: number;
  startItem: number;
  endItem: number;
};

const PaginationCount = ({
  totalItems,
  startItem,
  endItem
}: PaginationCountProps) => {
  const { t } = useTranslation(['common']);
  return (
    <p className="text-gray-600 typo-para-medium">
      {t('pagination-count', {
        start: startItem,
        end: endItem,
        total: totalItems
      })}
    </p>
  );
};

export default PaginationCount;
