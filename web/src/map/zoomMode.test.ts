import { describe, expect, it } from 'vitest';
import { DETAIL_MIN_ZOOM, NODE_CLUSTER_MAX_ZOOM, isClusterZoom, isDetailZoom, visualModeForZoom } from './zoomMode';

describe('map zoom visual modes', () => {
  it('uses one boundary for cluster and detail modes', () => {
    expect(NODE_CLUSTER_MAX_ZOOM).toBe(7);
    expect(DETAIL_MIN_ZOOM).toBeCloseTo(7.08);
    expect(visualModeForZoom(DETAIL_MIN_ZOOM - 0.01)).toBe('cluster');
    expect(visualModeForZoom(DETAIL_MIN_ZOOM)).toBe('detail');
    expect(isClusterZoom(6.9)).toBe(true);
    expect(isDetailZoom(8)).toBe(true);
  });
});
