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

export const formatFileSize = (size: number): string => {
  if (size === 0) return '0 Bytes';

  const units = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB'];
  const i = Math.floor(Math.log(size) / Math.log(1024));

  const formattedSize = (size / Math.pow(1024, i)).toFixed(2); // Two decimal places
  return `${formattedSize} ${units[i]}`;
};

export const convertFileToUnit8Array = (
  file: Blob,
  onLoad: (data: Uint8Array) => void
) => {
  const reader = new FileReader();
  reader.readAsArrayBuffer(file);
  reader.onload = () => {
    onLoad(new Uint8Array(reader.result as ArrayBuffer));
  };
};
