export class APIError extends Error {
  constructor(
    public status: number,
    message: string,
  ) {
    super(message);
  }
}
export async function request<T = any>(path: string, init: RequestInit = {}): Promise<T> {
  const headers = new Headers(init.headers);
  if (init.body && !(init.body instanceof FormData))
    headers.set('Content-Type', 'application/json');
  const res = await fetch('/api' + path, { ...init, headers, credentials: 'same-origin' });
  const payload = await res.json().catch(() => null);
  if (!res.ok)
    throw new APIError(res.status, payload?.error?.message || 'Permintaan gagal. Coba lagi.');
  return payload;
}
export function mutate<T = any>(path: string, body: any = {}, method = 'POST') {
  return request<T>(path, { method, body: JSON.stringify(body) });
}
export async function upload(file: File, purpose = 'media') {
  const form = new FormData();
  form.set('file', file);
  form.set('purpose', purpose);
  return request<{ id: string; url: string }>('/uploads', { method: 'POST', body: form });
}

export function recordClick(kind: 'ticket' | 'merch', target_id: string) {
  // Navigation does not wait for analytics; only confirmed API requests are counted.
  void mutate('/activity/click', { kind, target_id }).catch(() => {});
}
