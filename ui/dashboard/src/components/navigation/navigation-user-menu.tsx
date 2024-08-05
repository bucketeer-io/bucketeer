import {
  IconBuildingOutlined,
  IconLogoutOutlined,
  IconUserOutlined
} from '@icons';
import Menu from 'components/menu';

const NavigationUserMenu = () => {
  return (
    <div className="bg-primary-600">
      <Menu
        className="w-[200px]"
        options={[
          {
            text: 'User Profile',
            iconLeft: IconUserOutlined
          },
          {
            text: 'Polaris Edge',
            iconLeft: IconBuildingOutlined
            // iconRight: IconChevronDownOutlined
          },
          {
            text: 'Logout',
            iconLeft: IconLogoutOutlined
          }
        ]}
      />
    </div>
  );
};

export default NavigationUserMenu;
