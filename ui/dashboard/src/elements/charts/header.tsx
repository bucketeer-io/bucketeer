import { Tabs, TabsList, TabsTrigger } from 'components/tabs';
import ViewActions, { ViewActionsProps } from './view-actions';

export type Options = {
  label: string;
  value: string;
};

export type HeaderProps = ViewActionsProps & {
  title?: string;
  tabValue?: string;
  tabs?: Options[];
  onChangeTabs?: (value: string) => void;
};

const Header = ({
  title,
  tabValue,
  tabs,
  onChangeTabs,
  ...props
}: HeaderProps) => {
  return (
    <div className="flex items-center justify-between w-full p-5 gap-x-20">
      {title ? (
        <h3 className="typo-head-bold-big text-gray-900 whitespace-nowrap">
          {title}
        </h3>
      ) : (
        <Tabs value={tabValue} onValueChange={onChangeTabs}>
          <TabsList>
            {tabs?.map((item, index) => (
              <TabsTrigger key={index} value={item.value} className="pb-4">
                {item.label}
              </TabsTrigger>
            ))}
          </TabsList>
        </Tabs>
      )}
      <ViewActions {...props} />
    </div>
  );
};

export default Header;
