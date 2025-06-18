declare namespace ymaps {
  interface Map {
    destroy(): void;
    geoObjects: GeoObjectCollection;
    setCenter(center: number[], zoom: number): void;
  }

  interface GeoObjectCollection {
    add(object: Placemark): void;
  }

  interface MapConstructor {
    new (element: string | HTMLElement, options: MapOptions): Map;
  }

  interface MapOptions {
    center: number[];
    zoom: number;
    controls?: string[];
  }

  interface Placemark {
    events: {
      add(event: string, callback: () => void): void;
    };
  }

  interface PlacemarkConstructor {
    new (
      coordinates: number[],
      properties: { [key: string]: any },
      options: { [key: string]: any }
    ): Placemark;
  }

  interface TemplateLayoutFactory {
    createClass(template: string): any;
  }
}

interface Window {
  ymaps: {
    ready: (callback: () => void) => void;
    Map: ymaps.MapConstructor;
    Placemark: ymaps.PlacemarkConstructor;
    templateLayoutFactory: ymaps.TemplateLayoutFactory;
  };
} 