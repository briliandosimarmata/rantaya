import { api, query, requireUser } from '#lib/server/api.js';
import { error } from '@sveltejs/kit';
export const load = async (e: any) => {
  const { city } = await e.parent();
  const q = e.url.searchParams;
  const events = await api(
    e,
    '/events' +
      query({
        city,
        period: q.get('period') || 'upcoming',
        category: q.get('category') || '',
        from: q.get('from') || '',
        to: q.get('to') || '',
        page: q.get('page') || '1',
      }),
  );
  const [nearby, organizers] = await Promise.all([
    api(e, '/events' + query({ city, period: 'upcoming' })),
    api(e, '/organizers' + query({ city })),
  ]);
  return {
    events,
    nearby,
    organizers,
    period: q.get('period') || 'upcoming',
    category: q.get('category') || '',
    from: q.get('from') || '',
    to: q.get('to') || '',
    page: Number(q.get('page') || 1),
  };
};
