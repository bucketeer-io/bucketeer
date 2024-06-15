import { configureStore, getDefaultMiddleware } from '@reduxjs/toolkit';

import { thunkErrorHandler } from '../middlewares/thunkErrorHandler';
import { reducers } from '../modules';

export const store = configureStore({
  reducer: reducers,
  devTools: process.env.NODE_ENV === 'development',
  middleware: [...getDefaultMiddleware({}), thunkErrorHandler] as const,
});

export type AppDispatch = typeof store.dispatch;
