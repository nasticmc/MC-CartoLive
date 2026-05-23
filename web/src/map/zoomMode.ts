export const NODE_CLUSTER_MAX_ZOOM = 7;
export const DETAIL_MIN_ZOOM = NODE_CLUSTER_MAX_ZOOM;

export type MapVisualMode = 'cluster' | 'detail';

export function visualModeForZoom(zoom: number): MapVisualMode {
  return zoom < DETAIL_MIN_ZOOM ? 'cluster' : 'detail';
}

export function isClusterZoom(zoom: number): boolean {
  return visualModeForZoom(zoom) === 'cluster';
}

export function isDetailZoom(zoom: number): boolean {
  return visualModeForZoom(zoom) === 'detail';
}
