export type OrderBy =
  | 'DEFAULT'
  | 'ID'
  | 'CREATED_AT'
  | 'UPDATED_AT'
  | 'NAME'
  | 'URL_CODE';

export type OrderDirection = 'ASC' | 'DESC';

export type AddonSlot = 'left' | 'right';

// Theme
export type Color =
  | 'primary-600'
  | 'primary-500'
  | 'primary-300'
  | 'primary-200'
  | 'primary-100'
  | 'primary-50'
  | 'gray-600'
  | 'gray-500'
  | 'gray-300'
  | 'gray-200'
  | 'gray-100'
  | 'gray-50'
  | 'accent-red-600'
  | 'accent-red-500'
  | 'accent-red-300'
  | 'accent-red-200'
  | 'accent-red-100'
  | 'accent-red-50'
  | 'accent-orange-600'
  | 'accent-orange-500'
  | 'accent-orange-300'
  | 'accent-orange-200'
  | 'accent-orange-100'
  | 'accent-orange-50'
  | 'accent-green-600'
  | 'accent-green-500'
  | 'accent-green-300'
  | 'accent-green-200'
  | 'accent-green-100'
  | 'accent-green-50'
  | 'accent-blue-600'
  | 'accent-blue-500'
  | 'accent-blue-300'
  | 'accent-blue-200'
  | 'accent-blue-100'
  | 'accent-blue-50'
  | 'accent-pink-600'
  | 'accent-pink-500'
  | 'accent-pink-300'
  | 'accent-pink-200'
  | 'accent-pink-100'
  | 'accent-pink-50'
  | 'accent-yellow-600'
  | 'accent-yellow-500'
  | 'accent-yellow-300'
  | 'accent-yellow-200'
  | 'accent-yellow-100'
  | 'accent-yellow-50';

export type AvatarColor =
  | 'primary'
  | 'pink'
  | 'green'
  | 'blue'
  | 'orange'
  | 'red';

export type IconSize = 'xxs' | 'xs' | 'sm' | 'md' | 'lg' | 'xl' | '2xl' | '3xl';
