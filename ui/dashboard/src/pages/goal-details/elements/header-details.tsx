import { useToast } from 'hooks';
import { Goal } from '@types';
import { truncateTextCenter } from 'utils/converts';
import { copyToClipBoard } from 'utils/function';
import { cn } from 'utils/style';
import { IconCopy } from '@icons';
import Icon from 'components/icon';
import Status from 'elements/status';

const HeaderDetails = ({ goal }: { goal: Goal }) => {
  const { notify } = useToast();

  const handleCopyId = (id: string) => {
    copyToClipBoard(id);
    notify({
      toastType: 'toast',
      messageType: 'success',
      message: (
        <span>
          <b>ID</b> {` has been successfully copied!`}
        </span>
      )
    });
  };

  return (
    <div className="flex flex-col w-full gap-y-4 mt-3">
      <div className="flex items-center w-full gap-x-2">
        <h1 className="typo-head-bold-huge leading-6 text-gray-900">
          {goal.name}
        </h1>
        <Status
          status={goal?.isInUseStatus ? 'In Use' : 'Not In Use'}
          className={cn({
            'bg-accent-green-50 text-accent-green-500': goal.isInUseStatus
          })}
        />
      </div>
      <div className="flex items-center h-5 gap-x-2 typo-para-tiny text-gray-500 select-none">
        {truncateTextCenter(goal.id)}
        <div onClick={() => handleCopyId(goal.id)}>
          <Icon
            icon={IconCopy}
            size={'sm'}
            className="opacity-100 cursor-pointer"
          />
        </div>
      </div>
    </div>
  );
};

export default HeaderDetails;
