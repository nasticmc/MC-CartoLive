import type { CSSProperties } from 'react';
import { payloadLegendVisuals } from '../payloadVisuals';

export default function Legend() {
  const payloads = payloadLegendVisuals();
  return (
    <section className="legend-panel" aria-label="Map legend">
      <div className="legend-group">
        <span className="legend-title">Devices</span>
        <span><i className="legend-node repeater" />Repeaters</span>
        <span><i className="legend-node companion" />Companions</span>
        <span><i className="legend-node room" />Rooms</span>
        <span><i className="legend-node observer" />Observers</span>
      </div>
      <div className="legend-group">
        <span className="legend-title">Route color</span>
        <span className="frequency-ramp" />
        <span className="legend-scale"><b>Quiet</b><b>Busy</b></span>
      </div>
      <div className="legend-group packet-key">
        <span className="legend-title">Packet color</span>
        <div className="payload-key">
          {payloads.map((payload) => (
            <span className="payload-chip legend-payload" style={{ '--payload-color': payload.color } as CSSProperties} title={payload.description} key={payload.className}>
              <i />
              {payload.shortLabel}
            </span>
          ))}
        </div>
      </div>
    </section>
  );
}
