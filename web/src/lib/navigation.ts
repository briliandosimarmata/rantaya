// Keep a login/onboarding return destination on the same application origin.
export function returnPath(value: string | null | undefined, fallback = '/') {
  if (!value || !value.startsWith('/') || value.startsWith('//') || value.length > 2048) return fallback;
  try {
    const url = new URL(value, 'https://rantaya.invalid');
    if (url.origin !== 'https://rantaya.invalid' || /[\\\r\n]/.test(decodeURIComponent(url.pathname)) || decodeURIComponent(url.pathname).startsWith('//')) return fallback;
    return url.pathname + url.search + url.hash;
  } catch { return fallback; }
}
