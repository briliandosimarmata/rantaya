import type { OrderStatus } from './types';
export const money = (n: number) => n === 0 ? 'Gratis' : 'Rp' + Number(n || 0).toLocaleString('id-ID');
export const date = (
  s: string,
  options: Intl.DateTimeFormatOptions = { day: 'numeric', month: 'long', year: 'numeric' },
) =>
  s
    ? new Intl.DateTimeFormat('id-ID', { ...options, timeZone: 'Asia/Jakarta' }).format(new Date(s))
    : '';
export const time = (s: string) => date(s, { hour: '2-digit', minute: '2-digit' }) + ' WIB';
export const ended = (e: { ends_at: string }) => new Date(e.ends_at).getTime() < Date.now();
export const initials = (s: string) =>
  s
    .split(/\s+/)
    .slice(0, 2)
    .map((x) => x[0] || '')
    .join('')
    .toUpperCase();
export const status: Record<string, string> = {
  awaiting_payment: 'Menunggu pembayaran',
  awaiting_review: 'Menunggu konfirmasi',
  correction_requested: 'Perbaiki bukti transfer',
  approved: 'Tiket terbit',
  cancelled: 'Dibatalkan',
  expired: 'Waktu pembayaran berakhir',
};
export async function share(title: string, path: string) {
  const url = new URL(path, location.origin).href;
  if (navigator.share) await navigator.share({ title, url });
  else await navigator.clipboard.writeText(url);
}
export function video(url: string) {
  try {
    const u = new URL(url);
    const id =
      u.hostname === 'youtu.be'
        ? u.pathname.slice(1)
        : ['youtube.com', 'www.youtube.com', 'm.youtube.com'].includes(u.hostname)
          ? u.searchParams.get('v') || u.pathname.split('/').pop()
          : null;
    return id && /^[\w-]{11}$/.test(id)
      ? { type: 'youtube', url: 'https://www.youtube-nocookie.com/embed/' + id + '?autoplay=1' }
      : /\.(mp4|webm)$/i.test(u.pathname)
        ? { type: 'video', url }
        : null;
  } catch {
    return null;
  }
}
export function localDateValue(s: string) {
  const d = new Date(s);
  return new Intl.DateTimeFormat('sv-SE', {
    timeZone: 'Asia/Jakarta',
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  })
    .format(d)
    .replace(' ', 'T');
}

// Date-only filters must not shift with the browser's timezone.
export function dateRangeLabel(from: string, to: string, empty = 'Tanggal') {
  const short = (value: string) =>
    date(value + 'T12:00:00+07:00', { day: 'numeric', month: 'short' });
  return from && to
    ? from === to
      ? short(from)
      : short(from) + '–' + short(to)
    : from
      ? 'Mulai ' + short(from)
      : to
        ? 'Sampai ' + short(to)
        : empty;
}
