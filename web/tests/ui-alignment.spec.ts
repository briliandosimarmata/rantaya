import { test, expect } from '@playwright/test';

test('approved prototype layouts keep primary actions and mobile hierarchy', async ({
  page,
}, info) => {
  const mobile = info.project.name.startsWith('mobile');
  await page.goto('/agenda');
  await expect(page).toHaveTitle(/Rantaya$/);
  const poster = await page.locator('.event-poster').first().boundingBox();
  const content = await page.locator('.event-card-info').first().boundingBox();
  expect(poster!.x + poster!.width).toBeLessThan(content!.x);
  const filter = await page.locator('.agenda-date-sort').boundingBox();
  const chips = await page.locator('.filter-chips').boundingBox();
  expect(filter!.y).toBeLessThan(chips!.y + chips!.height);
  await page.screenshot({ path: info.outputPath('agenda.png'), fullPage: true });

  await page.goto('/event/di-balik-layar');
  const actions = page.locator(mobile ? '.mobile-ticket' : '.ticket-panel');
  const buy = actions.getByRole('link', { name: 'Beli tiket', exact: true });
  await expect(buy).toBeInViewport();
  if (mobile) {
    expect(await actions.evaluate((el) => getComputedStyle(el).position)).toBe('fixed');
    const nav = await page.locator('.bottom-nav').boundingBox();
    const bar = await actions.boundingBox();
    expect(bar!.y + bar!.height).toBeLessThanOrEqual(nav!.y + 1);
    const primary = await buy.boundingBox();
    const save = await actions
      .getByRole('button', { name: 'Simpan event', exact: true })
      .boundingBox();
    expect(Math.abs(primary!.y - save!.y)).toBeLessThan(2);
  }
  const organizerAvatar = await page.locator('.event-organizer-link .avatar').boundingBox();
  expect(Math.abs(organizerAvatar!.width - organizerAvatar!.height)).toBeLessThan(1);
  await page.screenshot({ path: info.outputPath('event.png'), fullPage: true });
  // Tablet has neither the desktop aside nor bottom navigation, but still exposes purchase.
  await page.setViewportSize({ width: 900, height: 900 });
  await expect(
    page.locator('.mobile-ticket').getByRole('link', { name: 'Beli tiket', exact: true }),
  ).toBeInViewport();
  await page.setViewportSize(info.project.use.viewport!);

  await page.goto('/ruang/teater-ruang');
  const header = await page.locator('.organizer-header').boundingBox();
  const description = await page.locator('.organizer-header .description').boundingBox();
  if (mobile) {
    expect(description!.width).toBeGreaterThan(header!.width - 2);
    const stats = await page.locator('.organizer-header .stats').boundingBox();
    const follow = await page.locator('.organizer-header .follow').boundingBox();
    expect(description!.y + description!.height).toBeLessThanOrEqual(stats!.y + 1);
    expect(stats!.y + stats!.height).toBeLessThanOrEqual(follow!.y + 1);
  }
  await page.screenshot({ path: info.outputPath('organizer.png'), fullPage: true });

  await page.goto('/tanggal?from=2027-03-05&to=2027-03-07');
  await expect(page.locator('.calendar-heading h2')).toHaveText(/Maret 2027/i);
  await page.getByRole('button', { name: 'Hari ini', exact: true }).click();
  await expect(page).toHaveURL(/\/tanggal\?/);
  await expect(page.locator('.date-selection')).not.toContainText(/\d{4}-\d{2}-\d{2}/);
  await expect(page.locator('.bottom-nav')).toHaveCount(0);
  await page.screenshot({ path: info.outputPath('calendar.png'), fullPage: true });
  await page.getByRole('button', { name: 'Tampilkan agenda', exact: true }).click();
  await expect(page).toHaveURL(/\/agenda\?/);
  const applied = page.url();
  await page.locator('.agenda-date-sort').click();
  await page.getByRole('button', { name: 'Bulan ini', exact: true }).click();
  await page.getByRole('button', { name: 'Kembali ke halaman sebelumnya' }).click();
  await expect(page).toHaveURL(applied);
});
