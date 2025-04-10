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

export const truncateBySide = (
  value: string,
  maxLen: number = 40,
  side: 'left' | 'right' = 'right'
) => {
  const valueLength = value.length;
  if (valueLength <= maxLen) return value;
  if (side === 'right') {
    const _value = value.slice(0, maxLen);
    return `${_value}...`;
  }
  const _value = value.slice(valueLength - maxLen, valueLength);
  return `...${_value}`;
};

export const formatFileSize = (size: number): string => {
  if (size === 0) return '0 Bytes';

  const units = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB'];
  const i = Math.floor(Math.log(size) / Math.log(1024));

  const formattedSize = (size / Math.pow(1024, i)).toFixed(2); // Two decimal places
  return `${formattedSize} ${units[i]}`;
};

export const uint8ArrayToBase64 = (uint8Array: Uint8Array): string => {
  const binary = Array.from(uint8Array)
    .map(byte => String.fromCharCode(byte))
    .join('');
  return btoa(binary);
};

export const covertFileToByteString = (
  file: Blob,
  onLoad: (data: Uint8Array) => void
) => {
  const reader = new FileReader();
  reader.readAsArrayBuffer(file);
  reader.onload = () => {
    onLoad(new Uint8Array(reader.result as ArrayBuffer));
  };
};

export const covertFileToUint8ToBase64 = (
  file: Blob,
  onLoad: (data: string) => void
) => {
  const reader = new FileReader();
  reader.readAsArrayBuffer(file);

  reader.onload = () => {
    const base64String = uint8ArrayToBase64(
      new Uint8Array(reader.result as ArrayBuffer)
    );
    onLoad(base64String);
  };
};

export const isJsonString = (str: string) => {
  try {
    const parsed = JSON.parse(str);
    if (parsed && typeof parsed === 'object') {
      return true;
    }
  } catch {
    return false;
  }
  return false;
};

export const areJsonStringsEqual = (json1: string, json2: string): boolean => {
  try {
    const obj1 = JSON.parse(json1);
    const obj2 = JSON.parse(json2);
    return JSON.stringify(obj1) === JSON.stringify(obj2);
  } catch {
    return false;
  }
};
