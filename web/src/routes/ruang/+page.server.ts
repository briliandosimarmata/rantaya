import { redirect } from '@sveltejs/kit';
import { api, query } from '#lib/server/api.js';
export const load = async (e: any) => {
  if (e.url.searchParams.get('city') === 'Karawang') {
    e.cookies.set('ruang_city','Karawang',{path:'/',httpOnly:false,sameSite:'lax',maxAge:31536000});
    redirect(303, '/ruang');
  }
  const { city, organizers } = await e.parent();
  const category = e.url.searchParams.get('category') || '';
  const nearby = await api(e, '/events' + query({ city, period: 'upcoming' }));
  return { organizers: organizers.filter((o: any) => !category || o.category === category), allOrganizers: organizers, nearby, category };
};
