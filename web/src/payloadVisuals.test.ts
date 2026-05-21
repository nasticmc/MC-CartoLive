import { describe, expect, it } from 'vitest';
import { hiddenPayloadCount, normalizePayloadType, payloadColor, payloadLegendVisuals, payloadVisual, payloadVisualsFor } from './payloadVisuals';

describe('payload visuals', () => {
  it('returns stable colors and labels for known payload types', () => {
    expect(payloadVisual('ADVERT')).toMatchObject({
      color: '#2dd4bf',
      label: 'Advert',
      shortLabel: 'ADV',
      className: 'payload-advert'
    });
    expect(payloadVisual('PLAIN_TEXT')).toMatchObject({
      color: '#38bdf8',
      label: 'Plain text',
      shortLabel: 'TXT',
      className: 'payload-plain-text'
    });
    expect(payloadColor('RETURNED_PATH')).toBe('#facc15');
  });

  it('normalizes common aliases into existing visual families', () => {
    expect(normalizePayloadType('text')).toBe('PLAIN_TEXT');
    expect(normalizePayloadType('trace-route')).toBe('TRACE');
    expect(normalizePayloadType('return_path')).toBe('RETURNED_PATH');
    expect(normalizePayloadType('command')).toBe('CONTROL');
  });

  it('gives unknown payload types deterministic fallback styling', () => {
    const first = payloadVisual('EXPERIMENTAL_PACKET');
    const second = payloadVisual('experimental packet');
    expect(first.color).toBe(second.color);
    expect(first.className).toBe('payload-experimental-packet');
    expect(first.label).toBe('Experimental Packet');
    expect(first.shortLabel).toBe('EP');
  });

  it('sorts and limits payload mixes for compact route rows', () => {
    const visuals = payloadVisualsFor(['GROUP_TEXT', 'ADVERT', 'ACK', 'PLAIN_TEXT'], 3);
    expect(visuals.map((visual) => visual.shortLabel)).toEqual(['ADV', 'TXT', 'GRP']);
    expect(hiddenPayloadCount(['GROUP_TEXT', 'ADVERT', 'ACK', 'PLAIN_TEXT'], visuals.length)).toBe(1);
  });

  it('provides a compact fixed legend set', () => {
    expect(payloadLegendVisuals().map((visual) => visual.shortLabel)).toEqual(['ADV', 'TXT', 'GRP', 'TRC', 'RET', 'ACK', 'CTL', 'OTH']);
  });
});
