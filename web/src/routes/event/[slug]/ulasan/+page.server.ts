import { api, query, requireUser } from '#lib/server/api.js';
import { error } from '@sveltejs/kit';
export const load = async (e: any) => {
  const user = await requireUser(e, 'customer');
  const event = await api(e, '/events/' + encodeURIComponent(e.params.slug));
  const reviews = await api(e, '/reviews' + query({ event: event.id }));
  return { event, review: reviews.find((r: any) => r.account_id === user.id) };
};
