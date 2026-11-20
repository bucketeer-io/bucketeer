import { IconPlus } from '@icons';
import Button from 'components/button';
import Icon from 'components/icon';

const CreateNewOptionButton = ({
  text,
  onClick
}: {
  text: string;
  onClick: () => void;
}) => (
  <Button
    type="button"
    variant="text"
    className="h-10 self-center w-full bg-white hover:bg-gray-100 sticky left-0 right-0 bottom-0 border-t border-gray-200"
    onClick={onClick}
  >
    <Icon icon={IconPlus} color="primary-500" size={'xs'} />
    {text}
  </Button>
);

export default CreateNewOptionButton;
