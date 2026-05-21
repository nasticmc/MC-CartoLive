import type { PublicNode, PublicRoute } from '../types';

interface Props {
  node: PublicNode | null;
  route: PublicRoute | null;
  connectedRoutes: PublicRoute[];
  onRouteSelect: (routeID: string) => void;
}

export default function SelectionDrawer({ node, route, connectedRoutes, onRouteSelect }: Props) {
  if (!node && !route) return null;
  return (
    <aside className="selection-drawer">
      {node && (
        <>
          <span className="eyebrow">{node.role.replaceAll('_', ' ')}</span>
          <h2>{node.label}</h2>
          <dl>
            <div><dt>Last seen</dt><dd>{formatRelative(node.lastSeen)}</dd></div>
            <div><dt>Activity</dt><dd>{node.activityCount.toLocaleString()} packets</dd></div>
            <div><dt>Regions</dt><dd>{node.iatasHeardIn.slice(0, 5).join(', ') || 'Unknown'}</dd></div>
          </dl>
          <div className="drawer-route-list">
            {connectedRoutes.slice(0, 6).map((item) => (
              <button type="button" key={item.id} onClick={() => onRouteSelect(item.id)}>
                <span className={`route-swatch bucket-${item.frequencyBucket}`} />
                <span>{item.from.nodeId === node.id ? item.to.label : item.from.label}</span>
                <em>{item.packetCount}</em>
              </button>
            ))}
          </div>
        </>
      )}
      {!node && route && (
        <>
          <span className="eyebrow">route</span>
          <h2>{route.from.label}{' -> '}{route.to.label}</h2>
          <dl>
            <div><dt>Packets</dt><dd>{route.packetCount.toLocaleString()}</dd></div>
            <div><dt>Distance</dt><dd>{route.distanceKm.toFixed(1)} km</dd></div>
            <div><dt>Last heard</dt><dd>{formatRelative(route.lastHeard)}</dd></div>
          </dl>
        </>
      )}
    </aside>
  );
}

function formatRelative(ms: number): string {
  const seconds = Math.max(0, Math.round((Date.now() - ms) / 1000));
  if (seconds < 60) return `${seconds}s ago`;
  const minutes = Math.round(seconds / 60);
  if (minutes < 60) return `${minutes}m ago`;
  const hours = Math.round(minutes / 60);
  if (hours < 24) return `${hours}h ago`;
  return new Intl.DateTimeFormat(undefined, { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' }).format(new Date(ms));
}
