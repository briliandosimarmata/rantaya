import { test, expect } from '@playwright/test';
import { readFileSync } from 'node:fs';
const reference = JSON.parse(
  readFileSync(new URL('./fixtures/mobile-home-reference.json', import.meta.url), 'utf8'),
);

async function login(page: import('@playwright/test').Page) {
  await page.goto('/masuk/customer');
  await page.getByRole('button', { name: 'Coba akun demo penonton' }).click();
  await page.waitForURL((url) => ['/', '/onboarding'].includes(url.pathname));
  if (new URL(page.url()).pathname === '/onboarding')
    await page.getByRole('button', { name: 'Lewati, lengkapi nanti' }).click();
  await expect(page).toHaveURL(/\/$/);
  await page.evaluate(() => document.fonts.ready);
}

test('mobile home matches measured SELA header, navigation and feed geometry', async ({
  page,
}, info) => {
  const width = page.viewportSize()!.width;
  test.skip(width > 760, 'Reference is the mobile prototype.');
  await login(page);
  const expected: any = (reference.viewports as any)[width];
  const selectors: Record<string, string> = {
    header: '.topbar-inner',
    logo: '.brand-mark',
    city: '.location-trigger',
    search: '.search-trigger',
    notification: '.notification-icon',
    profile: '.profile-access .avatar',
    heading: '.home-heading',
    title: 'h1',
    subtitle: '.home-heading p',
    tabs: '.home-feed .tabs',
    tab: '.home-feed .tabs a',
    filters: '.home-feed .filters',
    chip: '.home-feed .filter-chips a',
    sort: '.home-feed select',
    composer: '.compose-entry',
    composerAvatar: '.compose-entry .avatar',
    firstPost: '.post-card[data-post-id="p1"]',
    postHeader: '.post-card[data-post-id="p1"] .post-head',
    author: '.post-card[data-post-id="p1"] .post-author .author',
    postMeta: '.post-card[data-post-id="p1"] .post-author > .meta',
    postTitle: '.post-card[data-post-id="p1"] .post-content h2',
    postBody: '.post-card[data-post-id="p1"] .post-body',
    postImage: '.post-card[data-post-id="p1"] .post-image',
    footer: '.post-card[data-post-id="p1"] .post-actions',
    vote: '.post-card[data-post-id="p1"] .vote-action',
    nav: '.bottom-nav',
    navItem: '.bottom-nav a',
    fab: '.mobile-compose-fab',
  };
  const actual = await page.evaluate(
    (selectors) =>
      Object.fromEntries(
        Object.entries(selectors).map(([key, selector]) => {
          const el = document.querySelector(selector)!;
          const r = el.getBoundingClientRect(),
            s = getComputedStyle(el);
          return [
            key,
            {
              x: r.x,
              y: r.y,
              width: r.width,
              height: r.height,
              font: s.fontSize,
              lineHeight: s.lineHeight,
              padding: s.padding,
              radius: s.borderRadius,
              gap: s.gap,
              weight: s.fontWeight,
              background: s.backgroundColor,
            },
          ];
        }),
      ),
    selectors,
  );
  for (const key of [
    'header',
    'city',
    'search',
    'notification',
    'profile',
    'heading',
    'title',
    'subtitle',
    'tabs',
    'tab',
    'filters',
    'chip',
    'sort',
    'composer',
    'composerAvatar',
  ]) {
    for (const dimension of ['x', 'y', 'width', 'height'] as const)
      expect(
        Math.abs(actual[key][dimension] - expected[key][dimension]),
        key + '.' + dimension,
      ).toBeLessThan(0.1);
  }
  for (const dimension of ['x', 'width', 'height'] as const)
    expect(
      Math.abs(actual.postHeader[dimension] - expected.postHeader[dimension]),
      'postHeader.' + dimension,
    ).toBeLessThan(0.1);
  expect(actual.postHeader.y - actual.firstPost.y).toBeCloseTo(
    expected.postHeader.y - expected.firstPost.y,
    1,
  );
  for (const key of [
    'title',
    'subtitle',
    'tab',
    'chip',
    'sort',
    'composerAvatar',
    'author',
    'postMeta',
    'postTitle',
    'postBody',
    'vote',
    'navItem',
  ])
    expect(actual[key].font, key + ' font').toBe(expected[key].font);
  for (const key of ['city', 'composer', 'firstPost', 'postImage', 'navItem', 'fab'])
    expect(actual[key].radius, key + ' radius').toBe(expected[key].radius);
  expect(actual.firstPost.width).toBe(expected.firstPost.width);
  expect(actual.postImage.width / actual.postImage.height).toBeCloseTo(1.55, 2);
  expect(actual.nav.height).toBe(75);
  expect(actual.navItem.background).toBe('rgb(211, 233, 107)');
  expect(actual.fab.width).toBe(54);
  expect(
    await page
      .locator('.post-card[data-post-id="p1"] .date-tile')
      .evaluate((el) => el.getBoundingClientRect().height),
  ).toBe(55);
  expect(await page.locator('.mobile-brand').isVisible()).toBe(width > 370);
  expect(await page.locator('.mobile-discovery').count()).toBe(0);
  await expect(page.locator('.home-heading')).toContainText('Karawang · Ruang seni lokal');
  await expect(page.locator('.home-feed .tabs a')).toHaveText(['Jelajah', 'Diikuti']);
  await expect(page.locator('.composer-start')).toHaveText('Mau cerita apa hari ini?');
  await expect(page.locator('.post-card[data-post-id="p1"] .post-actions').first()).toContainText(
    'komentar',
  );
  const order = await page
    .locator('.post-card[data-post-id="p1"]')
    .evaluate((el) => [
      el.querySelector('.mention-chips')!.getBoundingClientRect().y,
      el.querySelector('.event-inline')!.getBoundingClientRect().y,
      el.querySelector('.post-card[data-post-id="p1"] .post-image')!.getBoundingClientRect().y,
    ]);
  expect(order[0]).toBeLessThan(order[1]);
  expect(order[1]).toBeLessThan(order[2]);
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true);
  await page.screenshot({ path: info.outputPath('home-aligned.png'), fullPage: true });
  await page.getByRole('link', { name: 'Diikuti', exact: true }).click();
  await expect(page).toHaveURL(/tab=following/);
  await expect(page.locator('.home-feed .tabs a.active')).toHaveText('Diikuti');
  await page.locator('.bottom-nav a[href="/agenda"]').click();
  await expect(page.locator('.bottom-nav a.active')).toHaveText('Agenda');
  await expect
    .poll(() =>
      page.locator('.bottom-nav a.active').evaluate((el) => getComputedStyle(el).backgroundColor),
    )
    .toBe('rgb(211, 233, 107)');
  expect(
    await page.locator('.bottom-nav a.active svg').evaluate((el) => getComputedStyle(el).boxShadow),
  ).toBe('none');
});

test('upvote reacts immediately, persists, reverses and rolls back on failure', async ({
  page,
}, info) => {
  const errors: string[] = [];
  page.on('pageerror', (error) => errors.push(error.message));
  await login(page);
  const card = page.locator('.post-card').first();
  const vote = card.getByRole('button', { name: 'Dukung postingan' });
  const original = Number(await vote.locator('span').textContent());
  await page.route('**/api/posts/*/vote', async (route) => {
    await new Promise((r) => setTimeout(r, 400));
    await route.continue();
  });
  await vote.click();
  await expect(vote).toHaveAttribute('aria-pressed', 'true');
  await expect(vote.locator('span')).toHaveText(String(original + 1));
  await expect(vote).toHaveAttribute('aria-busy', 'true');
  expect(await vote.locator('svg').evaluate((el) => el.getAnimations().length)).toBeGreaterThan(0);
  await expect
    .poll(() => vote.evaluate((el) => getComputedStyle(el).backgroundColor))
    .toBe('rgb(211, 233, 107)');
  await expect(vote).toHaveAttribute('aria-busy', 'false');
  await page.reload();
  await expect(vote).toHaveAttribute('aria-pressed', 'true');
  await expect(vote.locator('span')).toHaveText(String(original + 1));
  await page.screenshot({ path: info.outputPath('home-upvote-active.png'), fullPage: true });
  await vote.click();
  await expect(vote).toHaveAttribute('aria-busy', 'false');
  await expect(vote).toHaveAttribute('aria-pressed', 'false');
  await expect(vote.locator('span')).toHaveText(String(original));
  await page.unroute('**/api/posts/*/vote');
  await page.route('**/api/posts/*/vote', (route) =>
    route.fulfill({
      status: 503,
      contentType: 'application/json',
      body: JSON.stringify({ error: { message: 'Uji: coba kembali.' } }),
    }),
  );
  await vote.click();
  await expect(vote).toHaveAttribute('aria-busy', 'false');
  await expect(vote).toHaveAttribute('aria-pressed', 'false');
  await expect(vote.locator('span')).toHaveText(String(original));
  await expect(page.getByRole('status')).toContainText('Uji: coba kembali.');
  await page.unroute('**/api/posts/*/vote');
  await page.emulateMedia({ reducedMotion: 'reduce' });
  await vote.focus();
  await page.keyboard.press('Enter');
  await expect(vote).toHaveAttribute('aria-busy', 'false');
  await expect(vote).toHaveAttribute('aria-pressed', 'true');
  expect(await vote.locator('svg').evaluate((el) => getComputedStyle(el).animationName)).toBe(
    'none',
  );
  await vote.click();
  await expect(vote).toHaveAttribute('aria-busy', 'false');
  await expect(vote).toHaveAttribute('aria-pressed', 'false');
  expect(errors).toEqual([]);
});

test('home controls retain filtering, bookmarks, menu, search and global location', async ({
  page,
}) => {
  await login(page);
  await page.getByRole('button', { name: 'Pilihan postingan' }).first().click();
  await expect(page.locator('.post-menu')).toBeVisible();
  await expect(page.locator('.post-menu a')).toHaveText('Laporkan konten');
  await page.keyboard.press('Escape');
  await expect(page.locator('.post-menu')).toHaveCount(0);
  const save = page.locator('.post-card').first().getByRole('button', { name: 'Simpan postingan' });
  await save.click();
  await expect(save).toHaveAttribute('aria-pressed', 'true');
  await expect
    .poll(() => save.evaluate((el) => getComputedStyle(el).backgroundColor))
    .toBe('rgb(211, 233, 107)');
  await page.reload();
  await expect(save).toHaveAttribute('aria-pressed', 'true');
  await save.click();
  await expect(save).toHaveAttribute('aria-pressed', 'false');
  await page.getByLabel('Urutan postingan').selectOption('popular');
  await expect(page).toHaveURL(/sort=popular/);
  await page.locator('.home-feed .filter-chips a').filter({ hasText: 'Musik' }).click();
  await expect(page).toHaveURL(/category=Musik/);
  await expect(page.locator('.home-feed .filter-chips a.active')).toHaveText('Musik');
  await expect(page.locator('.post-card[data-post-id="p3"]')).toBeVisible();
  await expect(page.locator('.post-card[data-post-id="p2"]')).toHaveCount(0);
  await page.locator('.home-feed .filter-chips a').filter({hasText:'Teater'}).click();
  await expect(page.locator('.post-card[data-post-id="p2"]')).toBeVisible();
  await expect(page.locator('.post-card[data-post-id="p3"]')).toHaveCount(0);
  await page.locator('.home-feed .filter-chips a').filter({hasText:'Musik'}).click();
  await expect(page).toHaveURL(/category=Musik/);
  await expect(page.locator('.home-feed .filter-chips a.active')).toHaveText('Musik');
  await page.locator('.search-trigger').click();
  await expect(page).toHaveURL(/\/cari$/);
  await page.getByRole('button', { name: 'Kembali ke halaman sebelumnya' }).click();
  await expect(page).toHaveURL(/category=Musik/);
  await page.locator('.location-trigger').click();
  await expect(page.locator('#city-menu')).toContainText('Jelajahi kota');
  await page.locator('#city-menu button').filter({ hasText: 'Jakarta' }).click();
  await expect(page.locator('.home-heading')).toContainText('Jakarta · Ruang seni lokal');
  await expect(page.locator('.empty')).toContainText('Belum ada obrolan di sini');
  await page.reload();
  await expect(page.locator('.location-trigger')).toHaveText('Jakarta');
  await page.locator('.location-trigger').click();
  await page.locator('#city-menu button').filter({ hasText: 'Karawang' }).click();
  await expect(page.locator('.home-heading')).toContainText('Karawang · Ruang seni lokal');
});
