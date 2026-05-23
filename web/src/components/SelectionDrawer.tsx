import type { PublicNode, PublicRoute } from '../types';

interface Props {
  node: PublicNode | null;
  route: PublicRoute | null;
  connectedRoutes: PublicRoute[];
  allRoutes: PublicRoute[];
  onRouteSelect: (routeID: string) => void;
  onClose: () => void;
}

export default function SelectionDrawer({ node, route, connectedRoutes, allRoutes, onRouteSelect, onClose }: Props) {
  if (!node && !route) return null;
  return (
    <aside className="selection-drawer">
      <button className="drawer-close" type="button" aria-label="Close panel" onClick={onClose}>×</button>
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

          <ReachablePhonebook node={node} allRoutes={allRoutes} onRouteSelect={onRouteSelect} />
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


type ReachableEntry = { nodeID: string; label: string; hopCount: number; pathRouteIDs: string[] };

function ReachablePhonebook({ node, allRoutes, onRouteSelect }: { node: PublicNode; allRoutes: PublicRoute[]; onRouteSelect: (routeID: string) => void }) {
  const grouped = buildReachable(node, allRoutes);
  if (grouped.length === 0) return null;
  return (
    <section className="reachable-phonebook">
      <h3>Reachable nodes</h3>
      {grouped.map((group) => (
        <div key={group.hopCount} className="phonebook-group">
          <h4>{group.hopCount} hop{group.hopCount === 1 ? '' : 's'}</h4>
          {group.entries.map((entry) => (
            <button key={entry.nodeID} type="button" className="phonebook-row" onMouseEnter={() => entry.pathRouteIDs[0] && onRouteSelect(entry.pathRouteIDs[0])}>
              <span>{entry.label}</span>
              <em>{entry.pathRouteIDs.length} route{entry.pathRouteIDs.length === 1 ? '' : 's'}</em>
            </button>
          ))}
        </div>
      ))}
    </section>
  );
}

function buildReachable(node: PublicNode, routes: PublicRoute[]) {
  const adj = new Map<string, Array<{next: string; routeID: string}>>();
  for (const r of routes) {
    adj.set(r.from.nodeId, [...(adj.get(r.from.nodeId) ?? []), { next: r.to.nodeId, routeID: r.id }]);
    adj.set(r.to.nodeId, [...(adj.get(r.to.nodeId) ?? []), { next: r.from.nodeId, routeID: r.id }]);
  }
  const seen = new Set([node.id]);
  const queue: Array<{id: string; hops: number; pathRouteIDs: string[]}> = [{ id: node.id, hops: 0, pathRouteIDs: [] }];
  const entries: ReachableEntry[] = [];
  while (queue.length) {
    const current = queue.shift()!;
    for (const edge of adj.get(current.id) ?? []) {
      if (seen.has(edge.next)) continue;
      seen.add(edge.next);
      const nextPath = [...current.pathRouteIDs, edge.routeID];
      const hopCount = current.hops + 1;
      entries.push({ nodeID: edge.next, label: edge.next, hopCount, pathRouteIDs: nextPath });
      queue.push({ id: edge.next, hops: hopCount, pathRouteIDs: nextPath });
    }
  }
  const groups = Array.from(entries.reduce((map, e) => {
    const list = map.get(e.hopCount) ?? [];
    list.push(e);
    map.set(e.hopCount, list);
    return map;
  }, new Map<number, ReachableEntry[]>()));
  return groups.sort((a,b)=>a[0]-b[0]).map(([hopCount, entries]) => ({ hopCount, entries: entries.sort((a,b)=>a.label.localeCompare(b.label)).slice(0, 12) }));
}
