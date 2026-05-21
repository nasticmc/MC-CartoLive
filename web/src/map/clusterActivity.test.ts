import { describe, expect, it } from 'vitest';
import {
  CLUSTER_ACTIVITY_GLOW_MS,
  type ClusterActivityGlow,
  type ClusterActivityTarget,
  clusterActivityGlowID,
  clusterActivityGlowsToGeoJSON,
  clusterActivityIntensity,
  nearestClusterTarget,
  pruneClusterActivityGlows,
  upsertClusterActivityGlow
} from './clusterActivity';

const target = (clusterID: number, x: number, y: number): ClusterActivityTarget => ({
  clusterID,
  pointCount: 12,
  lng: -123.1 + clusterID,
  lat: 49.2,
  x,
  y
});

describe('cluster activity glows', () => {
  it('chooses the nearest rendered cluster target', () => {
    const nearest = nearestClusterTarget([target(1, 20, 20), target(2, 80, 80)], 72, 70);

    expect(nearest?.clusterID).toBe(2);
  });

  it('upserts payload-colored cluster glow state', () => {
    const glows = new Map<string, ClusterActivityGlow>();
    const id = upsertClusterActivityGlow(glows, target(9, 0, 0), 'GROUP_TEXT', 1000);

    expect(id).toBe(clusterActivityGlowID(9));
    expect(glows.get(id)).toMatchObject({
      color: '#a78bfa',
      expiresAt: 1000 + CLUSTER_ACTIVITY_GLOW_MS,
      pointCount: 12
    });
  });

  it('fades and prunes cluster glows without leaving stale features', () => {
    const glows = new Map<string, ClusterActivityGlow>();
    upsertClusterActivityGlow(glows, target(4, 0, 0), 'PLAIN_TEXT', 2000, 1000);

    const freshFeature = clusterActivityGlowsToGeoJSON(glows, 2000).features[0];
    const fadingFeature = clusterActivityGlowsToGeoJSON(glows, 2600).features[0];
    expect(freshFeature.properties.intensity).toBe(1);
    expect(fadingFeature.properties.intensity).toBeGreaterThan(0);
    expect(fadingFeature.properties.intensity).toBeLessThan(1);
    expect(clusterActivityIntensity(glows.values().next().value!, 3100)).toBe(0);
    expect(pruneClusterActivityGlows(glows, 3100)).toBe(0);
    expect(glows.size).toBe(0);
  });
});
