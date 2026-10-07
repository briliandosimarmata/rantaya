export const GET = () =>
  new Response(
    'User-agent: *\nDisallow: /api/\nDisallow: /kelola\nDisallow: /transaksi\nDisallow: /tiket\nDisallow: /profil\nDisallow: /admin\nDisallow: /tulis\nDisallow: /onboarding\nDisallow: /masuk\n',
    { headers: { 'content-type': 'text/plain; charset=utf-8' } },
  );
