import { base } from '#lib/server/api.js';
import type { RequestHandler } from '@sveltejs/kit';
const proxy: RequestHandler = async ({ request, params, url, fetch }) => {
  const headers = new Headers();
  for (const key of [
    'cookie',
    'authorization',
    'content-type',
    'accept',
    'origin',
    'sec-fetch-site',
    'idempotency-key',
  ]) {
    const value = request.headers.get(key);
    if (value) headers.set(key, value);
  }
  const upstream = await fetch(base() + '/api/' + params.path + url.search, {
    method: request.method,
    headers,
    body: ['GET', 'HEAD'].includes(request.method) ? undefined : await request.arrayBuffer(),
    redirect: 'manual',
  });
  const out = new Headers();
  for (const key of [
    'content-type',
    'cache-control',
    'location',
    'retry-after',
    'x-content-type-options',
  ]) {
    const value = upstream.headers.get(key);
    if (value) out.set(key, value);
  }
  for (const cookie of upstream.headers.getSetCookie()) out.append('set-cookie', cookie);
  return new Response(upstream.body, { status: upstream.status, headers: out });
};
export const GET = proxy;
export const POST = proxy;
export const PUT = proxy;
export const PATCH = proxy;
export const DELETE = proxy;
