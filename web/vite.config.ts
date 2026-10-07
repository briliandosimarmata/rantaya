import adapter from '@sveltejs/adapter-node';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
const origin = new URL(process.env.APP_URL || 'http://localhost:5173');
const trustedOrigins = (process.env.APP_ALLOWED_ORIGINS || '')
  .split(',')
  .map((value) => value.trim())
  .filter(Boolean);
for (const value of trustedOrigins) {
  const url = new URL(value);
  if (
    !['http:', 'https:'].includes(url.protocol) ||
    value !== url.origin ||
    url.hostname.includes('*')
  ) {
    throw new Error(
      'APP_ALLOWED_ORIGINS must contain exact HTTP(S) origins without credentials, paths or wildcards',
    );
  }
}
export default defineConfig({
  plugins: [
    sveltekit({
      adapter: adapter(),
      compilerOptions: { runes: true },
      // Kit 3 / adapter-node 6 embed this origin at build time. Keep CSRF on.
      paths: { origin: process.env.APP_URL || 'http://localhost:5173' },
      csrf: { trustedOrigins },
    }),
  ],
  server: {
    host: process.env.HOST || '127.0.0.1',
    port: Number(origin.port || (origin.protocol === 'https:' ? 443 : 80)),
    strictPort: true,
  },
});
