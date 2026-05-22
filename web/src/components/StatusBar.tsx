import type { CSSProperties } from 'react';
import { Database, MapPin, Route, Shield, Sparkles, Zap } from 'lucide-react';
import { payloadVisual } from '../payloadVisuals';
import type { LiveCoverageStats } from '../state';
import type { PublicStats } from '../types';
import { formatPacketsTotal, serverStatus } from './statusDisplay';

interface Props {
  stats: PublicStats | null;
  socketStatus: string;
  nodeCount: number;
  routeCount: number;
  coverage: LiveCoverageStats;
  latestPayloadTypeName: string | null;
  latestPacketID: string | null;
}

export default function StatusBar({ stats, socketStatus, nodeCount, routeCount, coverage, latestPayloadTypeName, latestPacketID }: Props) {
  const status = serverStatus(stats, socketStatus, coverage);
  const latestPayload = payloadVisual(latestPayloadTypeName);
  return (
    <header className="status-bar">
      <div className={`status-pill server-status ${status.live ? 'good' : 'warn'}`}>
        <span className={`server-signal ${status.live ? 'live' : 'stale'}`} />
        <span>{status.label}</span>
      </div>
      <div
        className="status-pill payload-signal-pill"
        style={{ '--payload-color': latestPayload.color } as CSSProperties}
        title={`Last packet type: ${latestPayload.label}`}
      >
        <span className="packet-type-signal" key={latestPacketID ?? latestPayloadTypeName ?? 'none'} />
        <span>{latestPayload.shortLabel}</span>
      </div>
      <div className="status-pill packets-total">
        <Database size={15} />
        <span>{formatPacketsTotal(stats?.packets)}</span>
      </div>
      <div className="status-pill pulse-rate">
        <Zap size={15} />
        <span>{coverage.receivedPerMinute}/min received</span>
      </div>
      <div className="status-pill route routed-rate">
        <Route size={15} />
        <span>{coverage.routeAnimatedPerMinute}/min routed</span>
      </div>
      <div className="status-pill observer">
        <Sparkles size={15} />
        <span>{coverage.observerBurstPerMinute}/min bursts</span>
      </div>
      <div className="status-pill unmapped">
        <MapPin size={15} />
        <span>{coverage.unmappedPerMinute}/min unresolved</span>
      </div>
      <div className="status-pill node-count">
        <Shield size={15} />
        <span>{nodeCount} nodes</span>
      </div>
      <div className="status-pill route route-count">
        <Route size={15} />
        <span>{routeCount} routes</span>
      </div>
    </header>
  );
}
