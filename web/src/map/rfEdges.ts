export { payloadColor } from '../payloadVisuals';

export function roleClass(role: string): string {
  if (role === 'repeater') return 'role-repeater';
  if (role === 'room_server') return 'role-room';
  if (role === 'companion') return 'role-companion';
  if (role === 'sensor') return 'role-sensor';
  return 'role-unknown';
}
