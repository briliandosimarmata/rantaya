import { api, query, requireUser } from '#lib/server/api.js';
export const load = async (e: any) => {
  await requireUser(e);
  const { city, organizers } = await e.parent();
  const [items, nearby] = await Promise.all([api(e, '/notifications'), api(e, '/events' + query({ city, period: 'upcoming' }))]);
  return { items, nearby, organizers };
};
