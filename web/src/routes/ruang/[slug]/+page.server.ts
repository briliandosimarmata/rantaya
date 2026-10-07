import { api, query, requireUser } from '#lib/server/api.js';
import { error } from '@sveltejs/kit';
export const load = async (e: any) => {
  const organizer = await api(e, '/organizers/' + encodeURIComponent(e.params.slug));
  const tab = e.url.searchParams.get('tab') || 'linimasa';
  const [posts, events, products, reviews] = await Promise.all([
    api(
      e,
      '/posts' +
        query({ organizer: organizer.id, mode: tab === 'mention' ? 'mentions' : 'official' }),
    ),
    api(e, '/events' + query({ organizer: organizer.id, period: 'all' })),
    api(e, '/products' + query({ organizer: organizer.id })),
    api(e, '/reviews' + query({ organizer: organizer.id })),
  ]);
  const { city } = await e.parent();
  const [nearby, organizers] = await Promise.all([
    api(e, '/events' + query({ city, period: 'upcoming' })),
    api(e, '/organizers' + query({ city })),
  ]);
  return { organizer, tab, posts, events, products, reviews, nearby, organizers };
};
