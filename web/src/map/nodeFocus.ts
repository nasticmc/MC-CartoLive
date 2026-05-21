import type { PublicRoute } from '../types';

export interface NodeFocus {
  selectedNodeID: string | null;
  neighbourNodeIDs: Set<string>;
  connectedRouteIDs: Set<string>;
  neighbourDistanceKmByNodeID: Map<string, number>;
}

export function nodeFocusFromRoutes(selectedNodeID: string | null, routes: PublicRoute[]): NodeFocus {
  const neighbourNodeIDs = new Set<string>();
  const connectedRouteIDs = new Set<string>();
  const neighbourDistanceKmByNodeID = new Map<string, number>();
  if (!selectedNodeID) {
    return { selectedNodeID: null, neighbourNodeIDs, connectedRouteIDs, neighbourDistanceKmByNodeID };
  }

  for (const route of routes) {
    const neighbourID = route.from.nodeId === selectedNodeID
      ? route.to.nodeId
      : route.to.nodeId === selectedNodeID
        ? route.from.nodeId
        : '';
    if (!neighbourID) continue;
    connectedRouteIDs.add(route.id);
    neighbourNodeIDs.add(neighbourID);
    const existingDistance = neighbourDistanceKmByNodeID.get(neighbourID);
    if (existingDistance === undefined || route.distanceKm < existingDistance) {
      neighbourDistanceKmByNodeID.set(neighbourID, route.distanceKm);
    }
  }

  return { selectedNodeID, neighbourNodeIDs, connectedRouteIDs, neighbourDistanceKmByNodeID };
}

export function emptyNodeFocus(): NodeFocus {
  return nodeFocusFromRoutes(null, []);
}
