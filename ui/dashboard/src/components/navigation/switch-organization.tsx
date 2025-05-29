import { useCallback, useEffect, useRef, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { switchOrganization } from '@api/auth';
import { useAuth } from 'auth';
import { PAGE_PATH_ROOT } from 'constants/routing';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import { clearCurrentEnvIdStorage } from 'storage/environment';
import { getOrgIdStorage, setOrgIdStorage } from 'storage/organization';
import { getTokenStorage, setTokenStorage } from 'storage/token';
import { cn } from 'utils/style';
import { IconChecked } from '@icons';
import Icon from 'components/icon';
import SearchInput from 'components/search-input';
import Spinner from 'components/spinner';

const OrganizationItem = ({
  name,
  active,
  isLoading,
  onClick
}: {
  name: string;
  active: boolean;
  isLoading: boolean;
  onClick: () => void;
}) => (
  <div
    className={cn(
      'flex items-center justify-between gap-x-2 px-3 py-2 text-gray-600 rounded-lg typo-para-medium cursor-pointer hover:bg-primary-400 hover:text-white',
      {
        'bg-primary-400 text-white': active,
        '!pointer-events-none': isLoading
      }
    )}
    onClick={() => !isLoading && onClick()}
  >
    <p className="line-clamp-1 break-all">{name}</p>
    {active &&
      (active && isLoading ? (
        <Spinner size="sm" className="min-w-5 size-5 border-2" />
      ) : (
        <Icon
          icon={IconChecked}
          size={'sm'}
          className="min-w-5 text-white flex-center"
        />
      ))}
  </div>
);

const SwitchOrganization = ({
  isOpen,
  onCloseSwitchOrg,
  onCloseSetting
}: {
  isOpen: boolean;
  onCloseSwitchOrg: () => void;
  onCloseSetting: () => void;
}) => {
  const navigate = useNavigate();
  const { t } = useTranslation(['common', 'form']);
  const { myOrganizations, onMeFetcher } = useAuth();
  const { errorNotify } = useToast();
  const organizationId = getOrgIdStorage();
  const [searchValue, setSearchValue] = useState('');
  const [organizations, setOrganizations] = useState(myOrganizations);
  const [currentOrganization, setCurrentOrganization] = useState(
    organizationId ?? ''
  );
  const [isLoading, setIsLoading] = useState(false);

  const menuRef = useRef<HTMLDivElement | null>(null);

  const onSearchOrganization = useCallback(
    (value: string) => {
      if (!value) return setOrganizations(myOrganizations);
      const newOrgs = myOrganizations.filter(item =>
        item.name?.toLowerCase()?.includes(value.toString())
      );
      setSearchValue(value);
      setOrganizations(newOrgs);
    },
    [myOrganizations]
  );

  const onChangeOrganization = useCallback(
    async (organizationId: string) => {
      try {
        setIsLoading(true);
        const token = getTokenStorage();
        if (token?.accessToken) {
          setOrgIdStorage(organizationId);
          clearCurrentEnvIdStorage();
          const resp = await switchOrganization({
            accessToken: token.accessToken,
            organizationId
          });
          if (resp.token) {
            setTokenStorage(resp.token);
            await onMeFetcher({ organizationId });
            onCloseSwitchOrg();
            onCloseSetting();
            navigate(PAGE_PATH_ROOT);
          }
        }
      } catch (error) {
        errorNotify(error);
      } finally {
        setIsLoading(false);
      }
    },
    [currentOrganization]
  );

  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (
        menuRef.current &&
        !menuRef.current.contains(event.target as Node) &&
        (event.target as Element).id !== 'switch-organization'
      ) {
        onCloseSwitchOrg();
      }
    }

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  useEffect(() => {
    if (!isOpen) setSearchValue('');
  }, [isOpen]);

  return (
    <div
      ref={menuRef}
      className={cn(
        'absolute z-50 top-0 left-[248px] w-[238px] h-screen bg-primary-100 transition-all duration-300',
        {
          'w-0 [&>div]:px-0 opacity-0': !isOpen
        }
      )}
    >
      <div
        className={cn(
          'flex flex-col size-full gap-y-5 p-4 overflow-y-auto relative small-scroll',
          {
            'overflow-hidden': isLoading
          }
        )}
      >
        <SearchInput
          variant="secondary"
          placeholder={`${t('form:placeholder-search')}`}
          value={searchValue}
          onChange={value => onSearchOrganization(value)}
        />
        {searchValue && !organizations?.length ? (
          <div className="flex flex-col justify-center items-center gap-3 pt-10 pb-4">
            <div className="typo-para-medium text-gray-500">
              {t(`navigation.no-organizations`)}
            </div>
          </div>
        ) : (
          <>
            <h3 className="typo-para-medium text-gray-600 whitespace-nowrap">
              {t('switch-organization')}
            </h3>
            <div className="flex flex-col gap-y-[1px]">
              {organizations?.map((item, index) => (
                <OrganizationItem
                  key={index}
                  name={item.name}
                  isLoading={isLoading}
                  active={currentOrganization === item.id}
                  onClick={() => {
                    if (currentOrganization === item.id) return;
                    onChangeOrganization(item.id);
                    setCurrentOrganization(item.id);
                  }}
                />
              ))}
            </div>
          </>
        )}
      </div>
    </div>
  );
};

export default SwitchOrganization;
