import { api, query } from '#lib/server/api.js';
export const load = async (event: any) => {
  const [account, config] = await Promise.all([api(event, '/auth/me'), api(event, '/auth/config')]);
  event.setHeaders({ 'cache-control': 'private, no-store' });
  const city = event.cookies.get('ruang_city') || 'Karawang';
  const allOrganizers = await api(event, '/organizers');
  const organizers = allOrganizers.filter((o: any) => o.city === city);
  const followedOrganizers = allOrganizers.filter((o: any) => account.follows.includes(o.id));
  return { account, config, city, organizers, followedOrganizers, origin: event.url.origin };
};
