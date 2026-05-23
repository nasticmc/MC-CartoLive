import { describe, expect, it } from 'vitest';
import { buildSharedViewURL, parseSharedView } from './shareView';

describe('share view URLs', () => {
  it('parses a viewport with selected route and search query', () => {
    expect(parseSharedView('?lat=43.6532&lng=-79.3832&z=9.5&route=r-1&q=Sydney')).toEqual({
      lat: 43.6532,
      lng: -79.3832,
      z: 9.5,
      route: 'r-1',
      node: undefined,
      q: 'Sydney'
    });
  });

  it('rejects invalid view coordinates', () => {
    expect(parseSharedView('?lat=200&lng=-79&z=9')).toBeNull();
    expect(parseSharedView('?lat=43&lng=-79&z=99')).toBeNull();
    expect(parseSharedView('?lng=-79&z=9')).toBeNull();
  });

  it('builds a share link with viewport, selected node, and query', () => {
    const url = buildSharedViewURL('https://routes.australiaverse.org/?old=1&route=old', { lat: 45.4215296, lng: -75.6971931, z: 8.125 }, {
      node: 'node-123',
      q: 'Melbourne'
    });
    expect(url).toBe('https://routes.australiaverse.org/?old=1&lat=45.42153&lng=-75.69719&z=8.13&node=node-123&q=Melbourne');
  });

  it('prefers route selection over node selection', () => {
    const url = buildSharedViewURL('https://routes.australiaverse.org/', { lat: 1, lng: 2, z: 3 }, {
      route: 'route-1',
      node: 'node-1'
    });
    expect(url).toBe('https://routes.australiaverse.org/?lat=1&lng=2&z=3&route=route-1');
  });
});
