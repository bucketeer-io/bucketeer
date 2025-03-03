export function formatBytes(bytes: number, decimals: number = 2) {
  if (!+bytes) return '0 Bytes';

  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return `${parseFloat((bytes / k ** i).toFixed(dm))} ${sizes[i]}`;
}

export const readImageFile = (file: File) => {
  return new Promise<string>(resolve => {
    const reader = new FileReader();
    reader.addEventListener(
      'load',
      () => resolve(reader.result as string),
      false
    );
    reader.readAsDataURL(file);
  });
};

export const convertBase64ToFile = async (
  base64Data: string,
  fileName: string = 'file'
) => {
  const base64Response = await fetch(base64Data);
  const blob = await base64Response.blob();
  const file = new File([blob], fileName, {
    type: blob.type
  });
  return file;
};
