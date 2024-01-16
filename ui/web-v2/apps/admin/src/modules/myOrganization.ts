import { Organization } from "@/proto/environment/organization_pb";
import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';

import { getMyOrganizations } from '../grpc/account';
import { GetMyOrganizationsRequest } from '../proto/account/service_pb';

const MODULE_NAME = 'myOrganization';

export const fetchMyOrganizations = createAsyncThunk<
  Array<Organization.AsObject>
>('me/fetchMyOrganizations', async () => {
  const request = new GetMyOrganizationsRequest();
  const res = await getMyOrganizations(request);
  return res.response.toObject().organizationsList;
});

export const myOrganizationSlice = createSlice({
  name: MODULE_NAME,
  initialState: {
    myOrganization: [],
  },
  reducers: {},
  extraReducers: (builder) => {
    builder.addCase(fetchMyOrganizations.fulfilled, (state, action) => {
      return {
        myOrganization: action.payload,
      };
    });
  },
});
