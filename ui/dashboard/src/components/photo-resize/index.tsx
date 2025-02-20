import { forwardRef, useCallback, useImperativeHandle, useState } from 'react';
import Cropper, { Area, Point } from 'react-easy-crop';
import { getCroppedImg } from './canvas-utils';

export interface PhotoResizeHandle {
  crop: () => void;
}

export interface PhotoResizeProps {
  aspect: number;
  value: string;
  cropShape?: 'round' | 'rect';
  onChange: (file: string) => void;
}

const PhotoResize = forwardRef<PhotoResizeHandle, PhotoResizeProps>(
  ({ aspect, cropShape = 'round', value, onChange }, ref) => {
    const [crop, setCrop] = useState<Point>({ x: 0, y: 0 });
    const [zoom, setZoom] = useState(1);
    const [croppedAreaPixels, setCroppedAreaPixels] = useState<Area>();

    const onCropComplete = useCallback(
      (_newCroppedArea: Area, newCroppedAreaPixels: Area) => {
        setCroppedAreaPixels(newCroppedAreaPixels);
      },
      []
    );

    const onConfirmCrop = useCallback(async () => {
      try {
        const croppedImage = await getCroppedImg(value, croppedAreaPixels!, 0);
        onChange(croppedImage!);
      } catch (e) {
        console.error(e);
      }
    }, [value, croppedAreaPixels]);

    useImperativeHandle(ref, () => ({
      crop: onConfirmCrop
    }));

    return (
      <div className="size-full">
        <div className="relative size-full">
          <Cropper
            image={value}
            crop={crop}
            zoom={zoom}
            aspect={aspect}
            cropShape={cropShape}
            showGrid={false}
            onCropChange={setCrop}
            onCropComplete={onCropComplete}
            onZoomChange={setZoom}
          />
        </div>
      </div>
    );
  }
);

export default PhotoResize;
