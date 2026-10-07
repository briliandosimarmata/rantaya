import { api } from '#lib/server/api.js';
const escape = (s: string) =>
  s.replace(
    /[&<>"']/g,
    (c) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&apos;' })[c]!,
  );
export const GET = async (e: any) => {
  let paths = ['/', '/agenda', '/ruang'];
  for (let page = 1; page <= 1000; page++) {
    const [events, organizers] = await Promise.all([
      api(e, '/events?period=all&page=' + page),
      api(e, '/organizers?page=' + page),
    ]);
    paths.push(
      ...events.filter((x: any) => x.published).map((x: any) => '/event/' + x.slug),
      ...organizers.map((x: any) => '/ruang/' + x.slug),
    );
    if (events.length < 20 && organizers.length < 20) break;
  }
  return new Response(
    '<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">' +
      paths.map((path) => '<url><loc>' + escape(e.url.origin + path) + '</loc></url>').join('') +
      '</urlset>',
    { headers: { 'content-type': 'application/xml', 'cache-control': 'public, max-age=300' } },
  );
};
