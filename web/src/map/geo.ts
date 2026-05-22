import type { PublicNode, PublicRouteEndpoint } from '../types';

export const AUSTRALIA_MAP_BOUNDS = {
  minLat: -44.5,
  maxLat: -9,
  minLng: 112,
  maxLng: 154
};

export function isMappableLatLng(lat: number, lng: number): boolean {
  return (
    Number.isFinite(lat) &&
    Number.isFinite(lng) &&
    lat !== 0 &&
    lng !== 0 &&
    lat >= AUSTRALIA_MAP_BOUNDS.minLat &&
    lat <= AUSTRALIA_MAP_BOUNDS.maxLat &&
    lng >= AUSTRALIA_MAP_BOUNDS.minLng &&
    lng <= AUSTRALIA_MAP_BOUNDS.maxLng
  );
}

export function isMappableNode(node: PublicNode): boolean {
  return isMappableLatLng(node.latitude, node.longitude);
}

export function isMappableEndpoint(endpoint: PublicRouteEndpoint): boolean {
  return isMappableLatLng(endpoint.lat, endpoint.lng);
}
