import { api, query, requireUser } from '#lib/server/api.js';
import { error } from '@sveltejs/kit';
export const load = async (e: any) => {
  const { user } = await e.parent();
  return {
    product: e.params.id === 'baru' ? null : await api(e, '/products/' + e.params.id),
    events: await api(e, '/events' + query({ organizer: user.organizer_id, period: 'all' })),
  };
};
