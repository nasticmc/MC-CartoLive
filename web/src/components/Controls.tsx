import { Eye, EyeOff, Pause, Play } from 'lucide-react';

interface Props {
  payloadTypes: string[];
  payloadFilters: Set<string>;
  setPayloadFilters: (next: Set<string>) => void;
  showUnresolved: boolean;
  setShowUnresolved: (next: boolean) => void;
  paused: boolean;
  setPaused: (next: boolean) => void;
}

export default function Controls({ payloadTypes, payloadFilters, setPayloadFilters, showUnresolved, setShowUnresolved, paused, setPaused }: Props) {
  const togglePayload = (name: string) => {
    const next = new Set(payloadFilters);
    if (next.has(name)) next.delete(name);
    else next.add(name);
    setPayloadFilters(next);
  };

  return (
    <div className="controls">
      <div className="control-row">
        <button className="control-button" type="button" onClick={() => setPaused(!paused)}>
          {paused ? <Play size={16} /> : <Pause size={16} />}
          <span>{paused ? 'Resume' : 'Pause'}</span>
        </button>
        <button className="control-button" type="button" onClick={() => setShowUnresolved(!showUnresolved)}>
          {showUnresolved ? <Eye size={16} /> : <EyeOff size={16} />}
          <span>{showUnresolved ? 'All' : 'High only'}</span>
        </button>
      </div>
      <div className="payload-grid">
        {payloadTypes.length === 0 && <span className="muted">Waiting for packets</span>}
        {payloadTypes.map((name) => (
          <label className="check-row" key={name}>
            <input type="checkbox" checked={payloadFilters.size === 0 || payloadFilters.has(name)} onChange={() => togglePayload(name)} />
            <span>{name.replaceAll('_', ' ')}</span>
          </label>
        ))}
      </div>
    </div>
  );
}
