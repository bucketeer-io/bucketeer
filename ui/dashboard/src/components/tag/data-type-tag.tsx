import { useMemo } from 'react';
import {
  IconDataObjectFilled,
  IconAbcFilled,
  IconToggleOnFilled
} from 'react-icons-material-design';
import { IconNumberOutlined } from '@icons';
import Icon from 'components/icon';

export type DataTypeTagType = 'toggle' | 'string' | 'number' | 'object';

export type DataTypeTagProps = {
  type: DataTypeTagType;
};

const DataTypeTag = ({ type }: DataTypeTagProps) => {
  const renderFlagIconType = useMemo(() => {
    switch (type) {
      case 'number':
        return IconNumberOutlined;
      case 'string':
        return IconAbcFilled;
      case 'toggle':
        return IconToggleOnFilled;
      case 'object':
        return IconDataObjectFilled;
    }
  }, [type]);

  return (
    <div className="bg-primary-50 size-8 rounded-[4px] grid place-items-center">
      <Icon icon={renderFlagIconType} size="sm" />
    </div>
  );
};

export default DataTypeTag;
