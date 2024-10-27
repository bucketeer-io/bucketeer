export type PaginationCountProps = {
  totalItems: number;
  value?: number;
};

const PaginationCount = ({ totalItems, value = 0 }: PaginationCountProps) => {
  return (
    <p className="text-gray-600 typo-para-medium">
      Showing <span>{value}</span> of <span>{totalItems}</span> results
    </p>
  );
};

export default PaginationCount;
