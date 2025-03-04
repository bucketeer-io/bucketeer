import { FC, useEffect, memo, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useLocation, useHistory } from 'react-router-dom';

import { AppState } from '../../modules';
import { exchangeTokenFromUrl } from '../../modules/auth';
import { AppDispatch } from '../../store';
import { getToken } from '../../storage/token';
import { fetchMyOrganizationsByAccessToken } from '../../modules/myOrganization';
import { parse } from 'query-string';
import { Organization } from '../../proto/environment/organization_pb';
import SelectOrganization from './selectOrganization';
import { setOrganizationId } from '../../storage/organizationId';
// import { fetchMe } from '../../modules/me';
// import { PAGE_PATH_ROOT } from '../../constants/routing';

export const AuthCallbackPage: FC = memo(() => {
  const history = useHistory();
  const dispatch = useDispatch<AppDispatch>();
  const location = useLocation();
  const loading = useSelector<AppState, boolean>((state) => state.auth.loading);
  const myOrganization = useSelector<AppState, Organization.AsObject[]>(
    (state) => state.myOrganization.myOrganization
  );
  const [selectedOrganization, setSelectedOrganization] = useState(null);

  useEffect(() => {
    const query = location.search;

    const { state, code } = parse(query);

    dispatch(
      fetchMyOrganizationsByAccessToken({ accessToken: code as string })
    ).then((res) => {
      const organizationList = res.payload as Organization.AsObject[];
      // if there is only one organization, set it as the default organization
      if (organizationList.length === 1) {
        setOrganizationId(organizationList[0].id);
        dispatch(
          exchangeTokenFromUrl({
            code: code as string,
            state: state as string,
            organizationId: organizationList[0].id
          })
        );
      }
    });
    // dispatch(exchangeTokenFromUrl(query));
  }, [dispatch]);

  useEffect(() => {
    const token = getToken();

    if (token?.accessToken) {
      history.push('/');
    }
  }, [loading]);

  const handleSubmit = () => {
    const query = location.search;
    const { state, code } = parse(query);

    setOrganizationId(selectedOrganization.value);
    dispatch(
      exchangeTokenFromUrl({
        code: code as string,
        state: state as string,
        organizationId: selectedOrganization.value
      })
    );
  };

  if (myOrganization.length > 1) {
    return (
      <SelectOrganization
        options={myOrganization.map((org) => ({
          label: org.name,
          value: org.id
        }))}
        onChange={(o) => setSelectedOrganization(o)}
        onSubmit={handleSubmit}
        isSubmitBtnDisabled={!selectedOrganization}
      />
    );
  }

  return <div className="spinner mx-auto mt-4" />;
});
