import { api, query, requireUser } from '#lib/server/api.js';
import { error } from '@sveltejs/kit';
export const load = async (e: any) => {
  await requireUser(e, 'customer');
  const event = await api(e, '/events/' + encodeURIComponent(e.params.slug));
  return {
    event,
    selectedSession: e.url.searchParams.get('session') || '',
    methods: await api(e, '/payment-methods' + query({ organizer: event.organizer_id })),
  };
};
