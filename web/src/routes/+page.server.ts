import { api, query } from '#lib/server/api.js';
import type { ArtsEvent, Post } from '#lib/types.js';
export const load = async (e: any) => {
  const { city } = await e.parent();
  const q = e.url.searchParams;
  const [posts, events, organizers] = await Promise.all([
    api(
      e,
      '/posts' +
        query({
          city,
          tab: q.get('tab') || '',
          sort: q.get('sort') || '',
          category: q.get('category') || '',
          page: q.get('page') || '1',
        }),
    ),
    api(e, '/events' + query({ city, period: 'upcoming' })),
    api(e, '/organizers' + query({ city })),
  ]);
  // Resolve event mentions outside the upcoming list without changing discovery pagination.
  const missing = [
    ...new Set(
      (posts as Post[]).flatMap((post) =>
        post.mentions.filter((mention) => mention.kind === 'event').map((mention) => mention.id),
      ),
    ),
  ].filter((id) => !(events as ArtsEvent[]).some((event) => event.id === id));
  const related = await Promise.all(missing.map((id) => api(e, '/events/' + id).catch(() => null)));
  return {
    posts,
    events,
    relatedEvents: [...events, ...related.filter(Boolean)],
    organizers,
    tab: q.get('tab') || 'explore',
    sort: q.get('sort') || 'latest',
    category: q.get('category') || '',
    page: Number(q.get('page') || 1),
  };
};
