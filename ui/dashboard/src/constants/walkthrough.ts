// Dispatched when the user cancels the leave-page confirmation dialog.
export const LEAVE_PAGE_CANCELLED_EVENT = 'leave-page-cancelled';

export const WALKTHROUGH_MOBILE_MENU_EVENT = 'walkthrough-mobile-menu';

export const WALKTHROUGH_TARGETS = {
  FEATURE_FLAGS_MENU: 'feature-flags-menu',
  CREATE_FLAG_BUTTON: 'create-flag-button',
  FLAG_GENERAL_INFO: 'flag-general-info',
  FLAG_VARIATIONS: 'flag-variations',
  SUBMIT_FLAG_BUTTON: 'submit-flag-button',
  CREATE_APIKEY_BUTTON: 'create-apikey-button',
  APIKEY_FORM: 'apikey-form',
  SUBMIT_APIKEY_BUTTON: 'submit-apikey-button'
} as const;
