import type { PublicLiveEnvelope } from './types';

export interface LiveSocket {
  close: () => void;
}

export function connectPublicSocket(onMessage: (message: PublicLiveEnvelope) => void, onStatus: (status: string) => void, onOpen?: () => void): LiveSocket {
  if (typeof window.WebSocket !== 'function') {
    onStatus('polling');
    return { close: () => undefined };
  }
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  const url = `${protocol}//${window.location.host}/ws/public`;
  let socket: WebSocket | null = null;
  let closed = false;
  let retryTimer: number | undefined;
  let attempts = 0;

  const connect = () => {
    if (closed) return;
    onStatus(attempts === 0 ? 'connecting' : 'retry');
    socket = new WebSocket(url);
    socket.addEventListener('open', () => {
      attempts = 0;
      onStatus('live');
      onOpen?.();
      socket?.send(JSON.stringify({ v: 1, type: 'subscribe', id: 'public-map' }));
    });
    socket.addEventListener('close', () => {
      if (closed) {
        onStatus('closed');
        return;
      }
      scheduleReconnect();
    });
    socket.addEventListener('error', () => onStatus('error'));
    socket.addEventListener('message', (event) => {
      try {
        onMessage(JSON.parse(event.data) as PublicLiveEnvelope);
      } catch {
        onStatus('bad-message');
      }
    });
  };

  const scheduleReconnect = () => {
    if (closed || retryTimer !== undefined) return;
    onStatus('retry');
    attempts += 1;
    const delay = Math.min(15_000, 800 * 2 ** Math.min(attempts, 5)) + Math.floor(Math.random() * 500);
    retryTimer = window.setTimeout(() => {
      retryTimer = undefined;
      connect();
    }, delay);
  };

  connect();

  return {
    close: () => {
      closed = true;
      if (retryTimer !== undefined) window.clearTimeout(retryTimer);
      socket?.close();
    }
  };
}
