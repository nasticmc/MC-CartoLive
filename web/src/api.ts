import type { PublicLiveState } from './types';

async function getJSON<T>(url: string): Promise<T> {
  const res = await fetch(url, { headers: { Accept: 'application/json' } });
  if (!res.ok) {
    throw new Error(`${res.status} ${res.statusText}`);
  }
  return res.json() as Promise<T>;
}

export function fetchPublicState(): Promise<PublicLiveState> {
  return getJSON<PublicLiveState>('/api/v1/public/state');
}
