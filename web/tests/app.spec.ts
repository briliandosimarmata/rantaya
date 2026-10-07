import { test, expect } from '@playwright/test';
import { fileURLToPath } from 'node:url';
const proofFile = fileURLToPath(new URL('./fixtures/proof.png', import.meta.url));
async function login(page: import('@playwright/test').Page, role: 'customer' | 'organizer') {
  await page.goto('/masuk/' + role);
  await page
    .getByRole('button', {
      name: role === 'customer' ? 'Coba akun demo penonton' : 'Coba akun demo pengelola',
    })
    .click();
  await page.waitForURL((u) => ['/', '/kelola', '/onboarding'].includes(u.pathname));
  if (new URL(page.url()).pathname === '/onboarding')
    await page.getByRole('button', { name: 'Lewati, lengkapi nanti' }).click();
  await expect(page).toHaveURL(role === 'customer' ? /\/$/ : /\/kelola$/);
}

test('organizer saves payment method, event, merch and profile through dashboard', async ({
  page,
}, testInfo) => {
  const errors: string[] = [];
  page.on('pageerror', (e) => errors.push(e.message));
  await login(page, 'organizer');
  const label = 'TEST ' + testInfo.project.name + ' ' + Date.now();
  await page.goto('/kelola/rekening');
  await page.locator('#provider').fill(label);
  await page.locator('#number').fill('0000000000');
  await page.locator('#holder').fill('TEST JANGAN TRANSFER');
  await page.getByRole('button', { name: 'Simpan metode' }).click();
  await expect(page.locator('article', { hasText: label })).toBeVisible();
  await page.reload();
  await expect(page.locator('article', { hasText: label })).toBeVisible();

  await page.goto('/kelola/event/baru');
  await page.locator('#event-title').fill(label);
  await page.locator('#event-description').fill('Pertunjukan sintetis untuk pengujian dashboard.');
  await page.locator('#event-venue').fill('Venue TEST');
  await page.locator('#event-maps').fill('https://maps.google.com/?q=Karawang');
  const start = new Date(Date.now() + 86400000).toISOString().slice(0, 16);
  const end = new Date(Date.now() + 90000000).toISOString().slice(0, 16);
  await page.locator('#session-start-0').fill(start);
  await page.locator('#session-end-0').fill(end);
  await page.getByRole('button', { name: 'Tambah pengisi', exact: true }).click();
  await page.locator('#performer-name-0').fill('Pemain TEST');
  await page.locator('#performer-role-0').fill('Aktor');
  await page.locator('#performer-photo-0').setInputFiles(proofFile);
  await expect(page.getByRole('img', { name: 'Pemain TEST', exact: true })).toBeVisible();
  await page
    .getByRole('checkbox', { name: 'Terbitkan event agar tampil di agenda publik.' })
    .check();
  await page.getByRole('button', { name: 'Simpan event', exact: true }).click();
  await expect(page).toHaveURL(/\/kelola\/event$/);
  const me = await (await page.request.get('/api/auth/me')).json();
  const events = await (
    await page.request.get('/api/events?period=all&organizer=' + me.user.organizer_id)
  ).json();
  const event = events.find((e: any) => e.title === label);
  expect(event).toBeTruthy();
  await page.goto('/event/' + event.slug);
  await expect(page.locator('h1')).toHaveText(label);
  await expect(page.getByText('Pemain TEST', { exact: true })).toBeVisible();
  await expect(page.locator('a[href="https://maps.google.com/?q=Karawang"]')).toBeVisible();
  await page.goto('/kelola/event/' + event.id);
  await expect(page.locator('#event-title')).toHaveValue(label);
  await page
    .locator('#event-description')
    .fill('Deskripsi yang diperbarui melalui dashboard pengelola.');
  await page.getByRole('button', { name: 'Simpan event', exact: true }).click();
  await expect(page).toHaveURL(/\/kelola\/event$/);
  await page.goto('/event/' + event.slug);
  await expect(
    page.getByText('Deskripsi yang diperbarui melalui dashboard pengelola.', { exact: true }),
  ).toBeVisible();

  await page.goto('/kelola/merch/baru');
  await page.locator('#product-name').fill(label);
  await page.locator('#product-desc').fill('Merchandise TEST untuk pengujian penyimpanan.');
  await page.locator('#product-price').fill('25000');
  await page.locator('#product-url').fill('https://example.com/merch');
  await page.locator('#product-event').selectOption(event.id);
  await page.getByRole('button', { name: 'Simpan merchandise' }).click();
  await expect(page).toHaveURL(/\/kelola\/merch$/);
  await expect(page.getByText(label, { exact: true })).toBeVisible();
  await page.reload();
  await expect(page.getByText(label, { exact: true })).toBeVisible();

  await page.goto('/kelola/profil');
  const about = 'Profil TEST tersimpan dari ' + label;
  await page.locator('#org-about').fill(about);
  await page.getByRole('button', { name: 'Simpan profil ruang' }).click();
  await expect(page.getByText('Profil ruang tersimpan.', { exact: true })).toBeVisible();
  await page.reload();
  await expect(page.locator('#org-about')).toHaveValue(about);
  expect(errors).toEqual([]);
});
test('public pages render semantic content within viewport', async ({ page }) => {
  const errors: string[] = [];
  page.on('pageerror', (e) => errors.push(e.message));
  for (const path of [
    '/',
    '/agenda',
    '/ruang',
    '/ruang/teater-ruang',
    '/event/di-balik-layar',
    '/cari',
    '/tanggal',
  ]) {
    const response = await page.goto(path);
    expect(response?.status()).toBe(200);
    await expect(page.locator('h1').first()).toBeVisible();
    expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth + 1)).toBe(
      true,
    );
  }
  expect(errors).toEqual([]);
});
test('customer composer has inline mentions and body-first focus', async ({ page }, testInfo) => {
  await page.goto('/masuk/customer');
  await page.getByRole('button', { name: 'Coba akun demo penonton' }).click();
  await page.waitForURL((u) => u.pathname === '/' || u.pathname === '/onboarding');
  if (new URL(page.url()).pathname === '/onboarding')
    await page.getByRole('button', { name: 'Lewati, lengkapi nanti' }).click();
  await expect(page).toHaveURL(/\/$/);
  await page.goto('/tulis');
  const body = page.locator('#post-body');
  await expect(body).toBeFocused();
  await expect(page.locator('.bottom-nav')).toHaveCount(0);
  const publish = await page.locator('.write-header').getByRole('button', {name:'Posting', exact:true}).boundingBox();
  const editor = await body.boundingBox();
  expect(publish!.y).toBeLessThan(editor!.y);
  await body.fill('Latihan malam ini @');
  await body.press('End');
  await body.dispatchEvent('input');
  await expect(page.locator('[role=option]').first()).toBeVisible();
  const box = await page.locator('#mention-results').boundingBox();
  expect(box?.y).toBeGreaterThan(0);
  expect((box?.y || 0) + (box?.height || 0)).toBeLessThanOrEqual(page.viewportSize()!.height + 1);
  await page.screenshot({path:testInfo.outputPath('composer.png'),fullPage:true});
  const option = page.locator('[role=option]').first();
  const name = await option.locator('strong').textContent();
  await option.click();
  expect(await body.inputValue()).toBe('Latihan malam ini @' + name + ' ');
  await expect(page.locator('#post-title')).toHaveCount(0);
  await body.fill('@tidak-ada-hasil-rantaya-xyz');
  await body.press('End'); await body.dispatchEvent('input');
  await expect(page.getByText('Tidak ada hasil. Coba nama lain.', {exact:true})).toBeVisible();
  await body.press('Escape');
  await expect(page.locator('#mention-results')).toHaveCount(0);
  // Simulate the reduced visualViewport of a software keyboard (not a physical keyboard test).
  await page.evaluate(() => Object.defineProperty(window, 'visualViewport', {configurable:true, value:{height:300,width:innerWidth,offsetTop:0,offsetLeft:0,addEventListener(){},removeEventListener(){}}}));
  await body.fill('Cerita latihan\n'.repeat(8) + '@');
  await body.press('End'); await body.dispatchEvent('input');
  await expect(page.locator('[role=option]').first()).toBeVisible();
  const keyboardBox = await page.locator('#mention-results').boundingBox();
  expect(keyboardBox!.y).toBeGreaterThanOrEqual(0);
  expect(keyboardBox!.y + keyboardBox!.height).toBeLessThanOrEqual(301);

});
test('organizer tools use separate account', async ({ page }) => {
  await page.goto('/masuk/organizer');
  await page.getByRole('button', { name: 'Coba akun demo pengelola' }).click();
  await page.waitForURL((u) => u.pathname === '/kelola' || u.pathname === '/onboarding');
  if (new URL(page.url()).pathname === '/onboarding')
    await page.getByRole('button', { name: 'Lewati, lengkapi nanti' }).click();
  await expect(page).toHaveURL(/\/kelola$/);
  for (const path of [
    '/kelola/event',
    '/kelola/event/baru',
    '/kelola/rekening',
    '/kelola/pembayaran',
    '/kelola/merch',
    '/kelola/merch/baru',
    '/kelola/ulasan',
    '/kelola/profil',
    '/kelola/scanner',
  ]) {
    expect((await page.goto(path))?.status()).toBe(200);
    expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth + 1)).toBe(
      true,
    );
  }
});

test('manual payment resumes, organizer approves and QR checks in once', async ({
  page,
  browser,
}) => {
  const errors: string[] = [];
  page.on('pageerror', (e) => errors.push(e.message));
  await page.goto('/masuk/customer');
  await page.waitForLoadState('networkidle');
  await page.getByRole('button', { name: 'Coba akun demo penonton' }).click();
  await page.waitForURL((u) => u.pathname === '/' || u.pathname === '/onboarding');
  if (new URL(page.url()).pathname === '/onboarding')
    await page.getByRole('button', { name: 'Lewati, lengkapi nanti' }).click();
  await expect(page).toHaveURL(/\/$/);
  await page.goto('/event/di-balik-layar/tiket');
  await page.locator('#quantity').fill('2');
  await page.getByRole('button', { name: 'Beli & lanjutkan pembayaran' }).click();
  await page.waitForURL(/\/transaksi\/[\w]+$/);
  const path = new URL(page.url()).pathname;
  const id = path.split('/').pop();
  await page.goto('/transaksi');
  await page.locator('a[href="' + path + '"]').click();
  await expect(page.locator('h2', { hasText: 'Bayar ke pengelola' })).toBeVisible();
  await page.locator('input[name=payment]').first().check();
  await expect(page.locator('#proof-file')).toBeEnabled();
  await page.locator('#proof-file').setInputFiles(proofFile);
  await page.getByRole('button', { name: 'Kirim bukti transfer' }).click();
  await expect(page.getByRole('heading', { name: 'Menunggu konfirmasi pengelola' })).toBeVisible();
  const organizerContext = await browser.newContext({
    baseURL: process.env.E2E_BASE_URL || 'http://localhost:5173',
    viewport: page.viewportSize()!,
  });
  const organizer = await organizerContext.newPage();
  organizer.on('pageerror', (e) => errors.push(e.message));
  try {
    await organizer.goto('/masuk/organizer');
    await organizer.waitForLoadState('networkidle');
    await organizer.getByRole('button', { name: 'Coba akun demo pengelola' }).click();
    await organizer.waitForURL((u) => u.pathname === '/kelola' || u.pathname === '/onboarding');
    if (new URL(organizer.url()).pathname === '/onboarding')
      await organizer.getByRole('button', { name: 'Lewati, lengkapi nanti' }).click();
    await expect(organizer).toHaveURL(/\/kelola$/);
    await organizer.goto('/kelola/pembayaran/' + id);
    const submitted = await (await page.request.get('/api/orders/' + id)).json();
    const proofPath = '/api/uploads/' + submitted.proof_id;
    const ownProof = await page.request.get(proofPath);
    expect(ownProof.status()).toBe(200);
    expect(ownProof.headers()['cache-control']).toContain('no-store');
    expect((await organizerContext.request.get(proofPath)).status()).toBe(200);
    const anonymous = await browser.newContext({
      baseURL: process.env.E2E_BASE_URL || 'http://localhost:5173',
    });
    try {
      expect((await anonymous.request.get(proofPath)).status()).toBe(401);
    } finally {
      await anonymous.close();
    }
    const orgNotifications = await (
      await organizerContext.request.get('/api/notifications')
    ).json();
    expect(
      orgNotifications.some(
        (n: any) => n.kind === 'payment' && n.url === '/kelola/pembayaran/' + id,
      ),
    ).toBe(true);

    await organizer.getByText('Minta perbaikan bukti', { exact: true }).click();
    await organizer.locator('#correction-note').fill('Bukti TEST kurang jelas, unggah ulang.');
    await organizer.getByRole('button', { name: 'Kirim permintaan perbaikan' }).click();
    await expect(organizer.getByText(/Menunggu penonton mengirim ulang bukti/)).toBeVisible();
    await page.reload();
    await expect(page.getByRole('heading', { name: 'Perbaiki bukti transfer' })).toBeVisible();
    await page.locator('#proof-file').setInputFiles(proofFile);
    await page.getByRole('button', { name: 'Kirim bukti transfer' }).click();
    await expect(
      page.getByRole('heading', { name: 'Menunggu konfirmasi pengelola' }),
    ).toBeVisible();
    await organizer.reload();
    await organizer.getByRole('checkbox', { name: /Saya sudah memeriksa mutasi/ }).check();
    await organizer.getByRole('button', { name: 'Setujui & terbitkan tiket' }).click();
    await expect(organizer.locator('article.ticket')).toHaveCount(2);
    await page.reload();
    await expect(page.locator('article.ticket')).toHaveCount(2);
    await expect(page.locator('article.ticket canvas').first()).toBeVisible();
    const userNotifications = await (await page.request.get('/api/notifications')).json();
    expect(userNotifications.filter((n: any) => n.url === path).length).toBe(5);
    const code = await page.locator('article.ticket code').first().textContent();
    const order = await (await organizerContext.request.get('/api/orders/' + id)).json();
    await organizer.goto('/kelola/scanner');
    await organizer.waitForLoadState('networkidle');
    await organizer.locator('#scan-event').selectOption(order.event.id);
    await organizer.locator('#ticket-code').fill(code!);
    await organizer.getByRole('button', { name: 'Validasi & catat masuk' }).click();
    await expect(organizer.getByText('Berhasil masuk', { exact: true })).toBeVisible();
    await organizer.locator('#ticket-code').fill(code!);
    await organizer.getByRole('button', { name: 'Validasi & catat masuk' }).click();
    await expect(organizer.getByText(/Tiket sudah digunakan pada/)).toBeVisible();
    expect(errors).toEqual([]);
  } finally {
    await organizerContext.close();
  }
});
