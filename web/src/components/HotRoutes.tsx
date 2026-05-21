import type { CSSProperties } from 'react';
import { hiddenPayloadCount, payloadVisualsFor } from '../payloadVisuals';
import type { RouteActivitySummary } from '../state';
import type { PublicRoute } from '../types';

interface Props {
  routes: PublicRoute[];
  selectedRouteID: string | null;
  routeActivityByID: Map<string, RouteActivitySummary>;
  onSelect: (routeID: string) => void;
}

export default function HotRoutes({ routes, selectedRouteID, routeActivityByID, onSelect }: Props) {
  return (
    <section className="hot-routes">
      <div className="panel-title compact">
        <span>Busy Pathways</span>
        <em>last 15m</em>
      </div>
      <div className="hot-route-list">
        {routes.slice(0, 10).map((route) => {
          const activity = routeActivityByID.get(route.id);
          const payloads = payloadVisualsFor(route.payloadTypeNames, 3);
          const hiddenCount = hiddenPayloadCount(route.payloadTypeNames, payloads.length);
          const hasRecentFlow = (activity?.total ?? 0) > 0;
          return (
            <button className={`hot-route ${route.id === selectedRouteID ? 'selected' : ''} ${hasRecentFlow ? 'flowing' : ''}`} key={route.id} type="button" onClick={() => onSelect(route.id)}>
              <span className={`route-swatch bucket-${route.frequencyBucket}`} />
              <span className="route-labels">
                <strong>{route.from.label}</strong>
                <span>{route.to.label}</span>
                <span className="hot-route-payloads" aria-label="Payload mix">
                  {payloads.map((payload) => (
                    <i className="payload-chip mini" style={{ '--payload-color': payload.color } as CSSProperties} title={payload.label} key={`${route.id}-${payload.className}`}>
                      {payload.shortLabel}
                    </i>
                  ))}
                  {hiddenCount > 0 && <i className="payload-chip mini muted-chip">+{hiddenCount}</i>}
                </span>
              </span>
              <RouteActivityBars activity={activity} />
              <em>{route.packetCount}</em>
            </button>
          );
        })}
        {routes.length === 0 && <div className="empty compact-empty">No busy pathways</div>}
      </div>
    </section>
  );
}

function RouteActivityBars({ activity }: { activity?: RouteActivitySummary }) {
  const bins = activity?.bins ?? Array.from({ length: 12 }, () => 0);
  const max = Math.max(1, ...bins);
  const total = activity?.total ?? 0;
  return (
    <span className="route-flow" title={`${total} packets in the last 15 minutes`}>
      <span className="route-spark">
        {bins.map((count, index) => (
          <i key={index} style={{ '--level': `${Math.max(0.08, count / max)}` } as CSSProperties} className={count > 0 ? 'active' : ''} />
        ))}
      </span>
      <span className={`flow-cue ${total > 0 ? 'visible' : ''}`}>{total > 0 ? 'Flow visible' : 'Quiet'}</span>
    </span>
  );
}
