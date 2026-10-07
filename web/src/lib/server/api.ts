import { error, redirect } from '@sveltejs/kit';
import type { RequestEvent } from '@sveltejs/kit';
export function base() {
  return process.env.API_BASE_URL || 'http://127.0.0.1:8080';
}
export async function api<T = any>(
  event: Pick<RequestEvent, 'fetch' | 'request'>,
  path: string,
): Promise<T> {
  let res: Response;
  try {
    res = await event.fetch(base() + '/api' + path, {
      headers: { cookie: event.request.headers.get('cookie') || '' },
    });
  } catch {
    error(503, 'Backend belum tersedia. Jalankan API Go dan PostgreSQL terlebih dahulu.');
  }
  const payload = await res.json().catch(() => null);
  if (!res.ok) error(res.status, payload?.error?.message || 'Data tidak dapat dimuat.');
  return payload;
}
export async function requireUser(event: any, role?: string) {
  const { account } = await event.parent();
  if (!account.user)
    redirect(
      303,
      '/masuk/' +
        (role === 'organizer' ? 'organizer' : 'customer') +
        '?next=' +
        encodeURIComponent(event.url.pathname + event.url.search),
    );
  if (role && account.user.role !== role) error(403, 'Halaman ini membutuhkan akun ' + role + '.');
  return account.user;
}
export const query = (o: Record<string, string | number | undefined>) =>
  '?' +
  new URLSearchParams(
    Object.entries(o)
      .filter(([, v]) => v !== undefined && v !== '')
      .map(([k, v]) => [k, String(v)]),
  );
