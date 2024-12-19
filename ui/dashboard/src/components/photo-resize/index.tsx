import { forwardRef, useCallback, useImperativeHandle, useState } from 'react';
import Cropper, { Area, Point } from 'react-easy-crop';
// import * as Slider from '@radix-ui/react-slider';
import { getCroppedImg } from './canvas-utils';

export interface PhotoResizeHandle {
  crop: () => void;
}

export interface PhotoResizeProps {
  aspect: number;
  value: string;
  onChange: (file: string) => void;
}

export type PhotoResizeRef = {};

const PhotoResize = forwardRef<PhotoResizeHandle, PhotoResizeProps>(
  ({ aspect, value, onChange }, ref) => {
    const [crop, setCrop] = useState<Point>({ x: 0, y: 0 });
    const [zoom, setZoom] = useState(1);
    const [croppedAreaPixels, setCroppedAreaPixels] = useState<Area>();

    const onSliderChangeZoom = (values: number[]) => {
      setZoom(values[0] || 1);
    };

    const onCropComplete = useCallback(
      (newCroppedArea: Area, newCroppedAreaPixels: Area) => {
        setCroppedAreaPixels(newCroppedAreaPixels);
      },
      []
    );
    const onConfirmCrop = useCallback(async () => {
      try {
        const croppedImage = await getCroppedImg(value, croppedAreaPixels, 0);
        onChange(croppedImage);
      } catch (e) {
        // eslint-disable-next-line no-console
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
            onCropChange={setCrop}
            onCropComplete={onCropComplete}
            onZoomChange={setZoom}
          />
        </div>

        {/* <Slider.Root
          className="relative mt-4 flex h-5 w-full touch-none select-none items-center justify-between"
          defaultValue={[1]}
          value={[zoom]}
          min={1}
          max={3}
          step={0.1}
          onValueChange={onSliderChangeZoom}
        >
          <Slider.Track className="relative h-1 w-full grow rounded-full bg-light-300">
            <Slider.Range className="absolute h-full rounded-full bg-light-600" />
          </Slider.Track>
          <Slider.Thumb className="block size-5 rounded-2xl border-2 border-primary-500 bg-light-50" />
        </Slider.Root> */}
      </div>
    );
  }
);

export default PhotoResize;
