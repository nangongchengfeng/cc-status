import { type ReactNode, useEffect, useRef, useState } from 'react';

interface ChartViewportSize {
  width: number;
  height: number;
}

interface ChartViewportProps {
  children: (size: ChartViewportSize) => ReactNode;
}

export function ChartViewport({ children }: ChartViewportProps) {
  const viewportRef = useRef<HTMLDivElement | null>(null);
  const [size, setSize] = useState<ChartViewportSize>({ width: 0, height: 0 });

  useEffect(() => {
    const element = viewportRef.current;

    if (!element) {
      return undefined;
    }

    const updateSize = () => {
      const nextWidth = Math.floor(element.getBoundingClientRect().width);
      const nextHeight = Math.floor(element.getBoundingClientRect().height);

      setSize((currentSize) => {
        if (currentSize.width === nextWidth && currentSize.height === nextHeight) {
          return currentSize;
        }

        return {
          width: nextWidth,
          height: nextHeight,
        };
      });
    };

    updateSize();

    if (typeof ResizeObserver === 'undefined') {
      return undefined;
    }

    const observer = new ResizeObserver(() => {
      updateSize();
    });

    observer.observe(element);

    return () => {
      observer.disconnect();
    };
  }, []);

  const isReady = size.width > 0 && size.height > 0;

  return <div ref={viewportRef} className="h-full w-full">{isReady ? children(size) : null}</div>;
}
