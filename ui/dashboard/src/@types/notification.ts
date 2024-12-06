export interface Notification {
  email: string;
  name: string;
  disabled: boolean;
  createdAt: string;
  lastName: string;
  language: string;
}

export interface NotificationsCollection {
  accounts: Array<Notification>;
  cursor: string;
  totalCount: string;
}
