import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';

import {
  getMyOrganizations,
  getMyOrganizationsByAccessToken
} from '../grpc/account';
import {
  GetMyOrganizationsRequest,
  GetMyOrganizationsByAccessTokenRequest
} from '../proto/account/service_pb';
import { Organization } from '../proto/environment/organization_pb';
import { AuthType } from '../proto/auth/service_pb';

const MODULE_NAME = 'myOrganization';

export const fetchMyOrganizations = createAsyncThunk<
  Array<Organization.AsObject>
>('me/fetchMyOrganizations', async () => {
  const request = new GetMyOrganizationsRequest();
  const res = await getMyOrganizations(request);
  return res.response.toObject().organizationsList;
});

interface FetchMyOrganizationsByAccessTokenParams {
  accessToken: string;
}

export const fetchMyOrganizationsByAccessToken = createAsyncThunk<
  Array<Organization.AsObject>,
  FetchMyOrganizationsByAccessTokenParams
>('me/fetchMyOrganizationsByAccessToken', async (params) => {
  const request = new GetMyOrganizationsByAccessTokenRequest();
  request.setAccessToken(params.accessToken);
  request.setType(AuthType.AUTH_TYPE_GOOGLE);

  const res = await getMyOrganizationsByAccessToken(request);
  return res.response.toObject().organizationsList;
});

export const myOrganizationSlice = createSlice({
  name: MODULE_NAME,
  initialState: {
    myOrganization: []
  },
  reducers: {},
  extraReducers: (builder) => {
    builder.addCase(fetchMyOrganizations.fulfilled, (state, action) => {
      return {
        myOrganization: action.payload
      };
    });
    builder.addCase(
      fetchMyOrganizationsByAccessToken.fulfilled,
      (state, action) => {
        return {
          myOrganization: action.payload
        };
      }
    );
  }
});
