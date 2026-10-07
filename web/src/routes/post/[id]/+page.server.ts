import { api, query, requireUser } from '#lib/server/api.js';
import { error } from '@sveltejs/kit';
export const load = async (e: any) => {
  const [post, comments] = await Promise.all([
    api(e, '/posts/' + e.params.id),
    api(e, '/posts/' + e.params.id + '/comments'),
  ]);
  const { city, organizers } = await e.parent();
  const nearby = await api(e, '/events' + query({ city, period: 'upcoming' }));
  const relatedEvents = await Promise.all(post.mentions.filter((m: any)=>m.kind==='event').map((m: any)=>api(e,'/events/'+m.id).catch((err: any)=>{if([403,404].includes(err.status))return null;throw err;})));
  return { post, comments, nearby, organizers, relatedEvents: relatedEvents.filter(Boolean) };
};
