export const joinName = (...args: (string | null | undefined)[]): string => {
  return args
    .filter(item => !!item)
    .join(' ')
    .trim();
};

export const getAvatarPlaceholder = (
  firstName: string | null,
  name?: string | null
) => {
  if (firstName && name) return `${firstName.trim()[0]}${name.trim()[0]}`;
  if (firstName) return firstName.trim()[0];
  if (name) return name.trim()[0];
  return '';
};
