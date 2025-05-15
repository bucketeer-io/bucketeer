import { cn } from 'utils/style';

const FeatureFlagStatus = ({
  status,
  enabled
}: {
  status: string;
  enabled: boolean;
}) => {
  return (
    <div
      className={cn(
        'flex-center py-0.5 px-2 rounded-lg typo-para-small !text-white !bg-primary-500 border border-gray-300',
        {
          '!text-gray-700 !bg-gray-200': !enabled
        }
      )}
    >
      {status}
    </div>
  );
};

export default FeatureFlagStatus;
