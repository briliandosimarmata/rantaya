import { test, expect } from '@playwright/test';

const uuidV4 = /^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/;

test.beforeEach(async ({ page }) => {
  await page.addInitScript(() => {
    Object.defineProperty(crypto, 'randomUUID', { configurable: true, value: undefined });
  });
  expect((await page.request.post('/api/auth/demo', { data: { role: 'customer' } })).status()).toBe(
    200,
  );
  await page.goto('/event/di-balik-layar/tiket');
  await page.waitForLoadState('networkidle');
});

test('checkout reaches payment without crypto.randomUUID and duplicate requests reuse the order', async ({
  page,
}) => {
  const errors: string[] = [];
  page.on('pageerror', (error) => errors.push(error.message));
  page.on('console', (message) => {
    if (message.type() === 'error') errors.push(message.text());
  });
  expect(await page.evaluate(() => typeof crypto.randomUUID)).toBe('undefined');
  const response = page.waitForResponse(
    (r) => new URL(r.url()).pathname === '/api/orders' && r.request().method() === 'POST',
  );
  await page.getByRole('button', { name: 'Beli & lanjutkan pembayaran' }).click();
  const created = await response;
  expect(created.status()).toBe(200);
  const payload = created.request().postDataJSON();
  expect(payload.idempotency_key).toMatch(uuidV4);
  const order = await created.json();
  await expect(page).toHaveURL(new RegExp(`/transaksi/${order.id}$`));
  await expect(page.getByRole('heading', { name: 'Bayar ke pengelola' })).toBeVisible();
  const duplicate = await page.request.post('/api/orders', { data: payload });
  expect(duplicate.status()).toBe(200);
  expect((await duplicate.json()).id).toBe(order.id);
  const orders = await (await page.request.get('/api/orders')).json();
  expect(orders.filter((entry: { id: string }) => entry.id === order.id)).toHaveLength(1);
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth + 1)).toBe(
    true,
  );
  expect(errors).toEqual([]);
});

test('retry after an ambiguous API failure keeps the key and recovers the original order', async ({
  page,
}) => {
  const errors: string[] = [];
  page.on('pageerror', (error) => errors.push(error.message));
  const keys: string[] = [];
  let originalOrder = '';
  await page.route('**/api/orders', async (route) => {
    if (route.request().method() !== 'POST') return route.continue();
    keys.push(route.request().postDataJSON().idempotency_key);
    if (keys.length === 1) {
      const upstream = await route.fetch();
      expect(upstream.status()).toBe(200);
      originalOrder = (await upstream.json()).id;
      await route.fulfill({
        status: 503,
        contentType: 'application/json',
        body: JSON.stringify({ error: { message: 'Koneksi TEST terputus. Coba lagi.' } }),
      });
    } else {
      await route.continue();
    }
  });
  const buy = page.getByRole('button', { name: 'Beli & lanjutkan pembayaran' });
  await buy.click();
  await expect(page.getByRole('alert')).toHaveText('Koneksi TEST terputus. Coba lagi.');
  await expect(buy).toBeEnabled();
  await buy.click();
  await expect(page).toHaveURL(/\/transaksi\/[\w]+$/);
  expect(new URL(page.url()).pathname).toBe('/transaksi/' + originalOrder);
  expect(keys).toHaveLength(2);
  expect(keys[0]).toMatch(uuidV4);
  expect(keys[1]).toBe(keys[0]);
  expect(errors).toEqual([]);
});

test('key generation failure shows an error and releases the busy button', async ({ page }) => {
  const errors: string[] = [];
  page.on('pageerror', (error) => errors.push(error.message));
  let orderRequests = 0;
  page.on('request', (request) => {
    if (request.method() === 'POST' && new URL(request.url()).pathname === '/api/orders')
      orderRequests++;
  });
  await page.evaluate(() => {
    Object.defineProperty(crypto, 'getRandomValues', {
      configurable: true,
      value: () => {
        throw new Error('Pembuat kunci TEST tidak tersedia.');
      },
    });
  });
  const buy = page.getByRole('button', { name: 'Beli & lanjutkan pembayaran' });
  await buy.click();
  await expect(page.getByRole('alert')).toHaveText('Pembuat kunci TEST tidak tersedia.');
  await expect(buy).toBeEnabled();
  expect(orderRequests).toBe(0);
  expect(errors).toEqual([]);
});
