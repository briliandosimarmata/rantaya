import { api, requireUser } from '#lib/server/api.js';
export const load = async (e: any) => {
  const user = await requireUser(e, 'organizer');
  return { user, organizer: await api(e, '/organizers/' + user.organizer_id) };
};
