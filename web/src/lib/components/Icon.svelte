<script lang="ts">
  let { name, size = 20 }: { name: string; size?: number } = $props();
  const paths: Record<string, string[]> = {
    home: ['m3 10 9-7 9 7v10H3z', 'M9 20v-7h6v7'],
    calendar: ['M3 5h18v16H3z', 'M7 3v4M17 3v4M3 11h18M7 15h3M14 15h3'],
    people: [
      'M9 4a4 4 0 1 0 0 8 4 4 0 0 0 0-8',
      'M2 21v-3a7 7 0 0 1 14 0v3M18 5a3 3 0 0 1 0 6M20 15v6',
    ],
    bookmark: ['M6 3h12v18l-6-4-6 4z'],
    bell: ['M5 17h14l-2-3V9a5 5 0 0 0-10 0v5zM10 21h4'],
    search: ['M10 4a6 6 0 1 0 0 12 6 6 0 0 0 0-12', 'm15 15 6 6'],
    pin: ['M20 10c0 6-8 11-8 11S4 16 4 10a8 8 0 1 1 16 0', 'M12 7a3 3 0 1 0 0 6 3 3 0 0 0 0-6'],
    back: ['m14 5-7 7 7 7'],
    next: ['m10 5 7 7-7 7'],
    down: ['m6 9 6 6 6-6'],
    up: ['m6 14 6-6 6 6'],
    chat: ['M21 11a8 8 0 0 1-8 8H7l-4 3V9a8 8 0 0 1 8-8h2a8 8 0 0 1 8 8z'],
    share: ['M12 15V3m-5 5 5-5 5 5M5 12v8h14v-8'],
    image: ['M3 3h18v18H3z', 'm3 17 5-5 4 4 4-6 5 8M7 7h2'],
    plus: ['M12 4v16M4 12h16'],
    edit: ['m4 15 12-12 5 5-12 12-6 1zM13 6l5 5'],
    flag: ['M5 22V3h14l-3 5 3 5H5'],
    close: ['m5 5 14 14M19 5 5 19'],
    check: ['m4 12 5 5L20 6'],
    ticket: ['M3 5h18v5a2 2 0 0 0 0 4v5H3v-5a2 2 0 0 0 0-4z', 'M15 5v3m0 3v2m0 3v3'],
    bag: ['M4 7h16l1 14H3zM8 7V5a4 4 0 0 1 8 0v2'],
    user: ['M12 3a4 4 0 1 0 0 8 4 4 0 0 0 0-8', 'M4 22v-3a8 8 0 0 1 16 0v3'],
    link: [
      'm9 15 6-6M8 17l-2 2a4 4 0 0 1-5-5l4-4a4 4 0 0 1 5 0M16 7l2-2a4 4 0 0 1 5 5l-4 4a4 4 0 0 1-5 0',
    ],
    scan: [
      'M3 8V3h5M16 3h5v5M21 16v5h-5M8 21H3v-5M7 7h3v3H7zM14 7h3v3h-3zM7 14h3v3H7zM14 14h3v3h-3z',
    ],
    settings: ['M4 6h16M4 12h16M4 18h16', 'M8 4v4M16 10v4M10 16v4'],
  };
</script>

<svg
  class="icon"
  width={size}
  height={size}
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width="1.65"
  stroke-linecap="round"
  stroke-linejoin="round"
  aria-hidden="true"
  >{#if name === 'calendar'}<rect x="3" y="5" width="18" height="16" rx="2" /><path
      d="M7 3v4M17 3v4M3 11h18M7 15h2M15 15h2"
    />
  {:else if name === 'people'}<circle cx="9" cy="8" r="3" /><path
      d="M3 20v-3a6 6 0 0 1 12 0v3M16 5a3 3 0 0 1 0 6M18 14a5 5 0 0 1 3 5"
    />
  {:else if name === 'search'}<circle cx="10.5" cy="10.5" r="6.5" /><path d="m16 16 5 5" />
  {:else if name === 'pin'}<path d="M20 10c0 6-8 11-8 11S4 16 4 10a8 8 0 1 1 16 0" /><circle
      cx="12"
      cy="10"
      r="2.5"
    />
  {:else if name === 'more'}<circle cx="5" cy="12" r="1" /><circle cx="12" cy="12" r="1" /><circle
      cx="19"
      cy="12"
      r="1"
    />
  {:else}{#each paths[name] || paths.chat as d}<path {d} />{/each}{/if}</svg
>
