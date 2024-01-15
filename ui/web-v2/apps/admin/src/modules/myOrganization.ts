import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';

import { getMyOrganizations } from '../grpc/account';
import { MyOrganization } from '../proto/account/account_pb';
import { GetMyOrganizationsRequest } from '../proto/account/service_pb';

const MODULE_NAME = 'myOrganization';

export const fetchMyOrganizations = createAsyncThunk<
  Array<MyOrganization.AsObject>
>('me/fetchMyOrganizations', async () => {
  const request = new GetMyOrganizationsRequest();
  const res = await getMyOrganizations(request);
  return res.response.toObject().myOrganizationsList;
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
