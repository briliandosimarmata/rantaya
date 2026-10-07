import { api, requireUser } from '#lib/server/api.js';
export const load = async (e: any) => {
  const user = await requireUser(e);
  const organizer = user.organizer_id ? await api(e, '/organizers/' + user.organizer_id) : null;
  return { tab: e.url.searchParams.get('tab') || 'profile', organizer };
};
