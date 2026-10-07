import { api, query, requireUser } from '#lib/server/api.js';
import { error } from '@sveltejs/kit';
export const load = async (e: any) => {
  const event = await api(e, '/events/' + encodeURIComponent(e.params.slug));
  const [posts, reviews] = await Promise.all([
    api(e, '/posts' + query({ event: event.id })),
    api(e, '/reviews' + query({ event: event.id })),
  ]);
  return { event, posts, reviews, tab: e.url.searchParams.get('tab') || 'tentang' };
};
