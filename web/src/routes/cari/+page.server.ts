import { api, query, requireUser } from '#lib/server/api.js';
import { error } from '@sveltejs/kit';
export const load = async (e: any) => {
  const { city } = await e.parent();
  const q = (e.url.searchParams.get('q') || '').slice(0, 100);
  if (!q.trim()) return { q, events: [], organizers: await api(e, '/organizers' + query({ city })), posts: [] };
  const [events, organizers, posts] = await Promise.all([
    api(e, '/events' + query({ city, q, period: 'all' })),
    api(e, '/organizers' + query({ city, q })),
    api(e, '/posts' + query({ city, q })),
  ]);
  return { q, events, organizers, posts };
};
