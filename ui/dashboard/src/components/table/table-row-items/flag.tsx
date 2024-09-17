import { IconManageAccountsOutlined } from 'react-icons-material-design';
import { Button } from 'components/button';
import Icon from 'components/icon';
import { DataTypeTag, StatusTag } from 'components/tag';
import { DataTypeTagType } from 'components/tag/data-type-tag';
import { StatusTagType } from 'components/tag/status-tag';
import Title, { TitleProps } from './title';

export type FlagProps = TitleProps &
  SubProps & {
    flagIconType?: DataTypeTagType;
    editable?: boolean;
  };

export type SubProps = {
  status?: StatusTagType;
  editable?: boolean;
};

const Flag = ({
  text,
  description,
  status,
  flagIconType,
  editable
}: FlagProps) => {
  return (
    <div className="flex gap-2">
      {flagIconType && <DataTypeTag type={flagIconType} />}
      <Title
        text={text}
        description={description}
        {...(status
          ? { sub: <Sub status={status} editable={editable} /> }
          : {})}
      />
    </div>
  );
};

export const Sub = ({ status, editable }: SubProps) => {
  return (
    <div className="flex items-center gap-2">
      {editable && (
        <Button
          variant="grey"
          size="icon-sm"
          className="size-[26px] bg-primary-50"
        >
          <Icon icon={IconManageAccountsOutlined} size="xs" />
        </Button>
      )}
      {status && <StatusTag variant={status} />}
    </div>
  );
};

export default Flag;
