import { IconInfoFilled } from '@icons';
import Icon from 'components/icon';

const InfoMessage = ({ description }: { description: string }) => {
  return (
    <div className="flex items-center w-full p-4 gap-x-2 rounded border-l-4 border-accent-blue-500 bg-accent-blue-50">
      <Icon icon={IconInfoFilled} size={'xxs'} color="accent-blue-500" />
      <p className="typo-para-small leading-[14px] text-accent-blue-500">
        {description}
      </p>
    </div>
  );
};

export default InfoMessage;
