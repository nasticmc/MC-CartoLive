export interface SharedViewState {
  lat: number;
  lng: number;
  z: number;
  route?: string;
  node?: string;
  q?: string;
}

export interface MapViewState {
  lat: number;
  lng: number;
  z: number;
}

export function parseSharedView(search: string): SharedViewState | null {
  const params = new URLSearchParams(search.startsWith('?') ? search : `?${search}`);
  if (!params.has('lat') || !params.has('lng') || !params.has('z')) return null;
  const lat = Number(params.get('lat'));
  const lng = Number(params.get('lng'));
  const z = Number(params.get('z'));
  if (!Number.isFinite(lat) || !Number.isFinite(lng) || !Number.isFinite(z)) return null;
  if (lat < -90 || lat > 90 || lng < -180 || lng > 180 || z < 0 || z > 24) return null;
  const route = params.get('route')?.trim() || undefined;
  const node = route ? undefined : params.get('node')?.trim() || undefined;
  const q = params.get('q')?.trim() || undefined;
  return { lat, lng, z, route, node, q };
}

export function buildSharedViewURL(baseHref: string, view: MapViewState, options: { route?: string | null; node?: string | null; q?: string }): string {
  const url = new URL(baseHref);
  url.searchParams.set('lat', fixedCoordinate(view.lat));
  url.searchParams.set('lng', fixedCoordinate(view.lng));
  url.searchParams.set('z', fixedZoom(view.z));
  url.searchParams.delete('route');
  url.searchParams.delete('node');
  if (options.route) {
    url.searchParams.set('route', options.route);
  } else if (options.node) {
    url.searchParams.set('node', options.node);
  }
  if (options.q?.trim()) url.searchParams.set('q', options.q.trim());
  else url.searchParams.delete('q');
  return url.toString();
}

function fixedCoordinate(value: number): string {
  return value.toFixed(5).replace(/\.?0+$/, '');
}

function fixedZoom(value: number): string {
  return value.toFixed(2).replace(/\.?0+$/, '');
}
