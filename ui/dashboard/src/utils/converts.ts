export const onGenerateSlug = (value: string) => {
  return value
    .toLowerCase()
    .trim()
    .replace(/[^\w\s-]/g, '')
    .replace(/\s+/g, '-')
    .replace(/-+/g, '-')
    .replace(/^-+|-+$/g, '');
};

export const onValidSlug = (value: string) => {
  return value
    .toLowerCase()
    .replace(/[^a-z0-9-]+/g, '')
    .replace(/-+/g, '-')
    .replace(/^-|-$/g, '');
};

export const truncateTextCenter = (value: string, maxLen: number = 14) => {
  if (value.length <= maxLen) return value;

  const half = Math.floor(maxLen / 2);
  const start = value.slice(0, half);
  const end = value.slice(value.length - half);

  return `${start}...${end}`;
};
