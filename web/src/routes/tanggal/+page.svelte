<script lang="ts">
  import { untrack } from 'svelte';
  import { goto } from '$app/navigation';
  import { date, dateRangeLabel } from '#lib/format.js';
  import SEO from '#lib/components/SEO.svelte';
  import Icon from '#lib/components/Icon.svelte';
  let { data } = $props();
  let from = $state(untrack(() => data.from)),
    to = $state(untrack(() => data.to));
  const now = new Date();
  const today = new Intl.DateTimeFormat('sv-SE', {
    timeZone: 'Asia/Jakarta',
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  }).format(now);
  const initial = untrack(() => (data.from || today).split('-').map(Number));
  let year = $state(initial[0]),
    month = $state(initial[1] - 1),
    presetChoice = $state(untrack(() => (data.from ? '' : 'all')));
  const key = (y: number, m: number, d: number) =>
    y + '-' + String(m + 1).padStart(2, '0') + '-' + String(d).padStart(2, '0');
  let cells = $derived.by(() => {
    const offset = (new Date(year, month, 1).getDay() + 6) % 7;
    const total = new Date(year, month + 1, 0).getDate();
    return Array.from({ length: offset + total }, (_, i) => (i < offset ? 0 : i - offset + 1));
  });
  function shift(n: number) {
    const d = new Date(year, month + n, 1);
    year = d.getFullYear();
    month = d.getMonth();
  }
  function pick(d: number) {
    presetChoice = '';
    const value = key(year, month, d);
    if (!from || to) {
      from = value;
      to = '';
    } else if (value < from) {
      to = from;
      from = value;
    } else to = value;
  }
  function apply() {
    goto(
      '/agenda?' +
        new URLSearchParams({ period: data.period, category: data.category, from, to: to || from }),
    );
  }
  function preset(kind: string) {
    presetChoice = kind;
    const [y, m, day] = today.split('-').map(Number);
    year = y;
    month = m - 1;
    if (kind === 'all') {
      from = '';
      to = '';
      return;
    }
    const start = new Date(y, m - 1, day, 12),
      end = new Date(start);
    if (kind === 'week') {
      start.setDate(start.getDate() - ((start.getDay() + 6) % 7));
      end.setTime(start.getTime());
      end.setDate(start.getDate() + 6);
    }
    if (kind === 'month') {
      start.setDate(1);
      end.setMonth(m, 0);
    }
    from = key(start.getFullYear(), start.getMonth(), start.getDate());
    to = key(end.getFullYear(), end.getMonth(), end.getDate());
  }
</script>

<SEO title="Pilih tanggal agenda" privatePage />
<section class="date-picker-page">
  <h1>Pilih tanggal</h1>
  <div class="date-presets">
    {#each [['all', 'Semua tanggal'], ['today', 'Hari ini'], ['week', 'Pekan ini'], ['month', 'Bulan ini']] as [value, label]}<button
        class:active={presetChoice === value}
        aria-pressed={presetChoice === value}
        onclick={() => preset(value)}>{label}</button
      >{/each}
  </div>
  <div class="date-selection" role="status" aria-live="polite">
    <span class="eyebrow">{from && !to ? 'Tanggal mulai dipilih' : 'Pilihan tanggal'}</span>
    <strong>{dateRangeLabel(from, to || from, 'Semua tanggal')}</strong><span class="meta"
      >{from && !to
        ? 'Pilih tanggal akhir, atau terapkan untuk satu hari.'
        : 'Pilih satu tanggal atau rentang tanggal.'}</span
    >
  </div>
  <div class="calendar">
    <div class="calendar-heading">
      <button class="icon-button" onclick={() => shift(-1)} aria-label="Bulan sebelumnya"
        ><Icon name="back" /></button
      >
      <h2>
        {new Intl.DateTimeFormat('id-ID', { month: 'long', year: 'numeric' }).format(
          new Date(year, month, 1),
        )}
      </h2>
      <button class="icon-button" onclick={() => shift(1)} aria-label="Bulan berikutnya"
        ><Icon name="next" /></button
      >
    </div>
    <div class="calendar-grid" role="group" aria-label="Kalender tanggal pertunjukan">
      {#each ['Sen', 'Sel', 'Rab', 'Kam', 'Jum', 'Sab', 'Min'] as day}<span>{day}</span>{/each}
      {#each cells as d}{#if d}<button
            class:selected={key(year, month, d) === from || key(year, month, d) === to}
            class:in-range={!!to && key(year, month, d) > from && key(year, month, d) < to}
            class:is-today={key(year, month, d) === today}
            onclick={() => pick(d)}
            aria-pressed={key(year, month, d) === from || key(year, month, d) === to}
            aria-current={key(year, month, d) === today ? 'date' : undefined}
            aria-label={date(key(year, month, d) + 'T12:00:00+07:00')}>{d}</button
          >{:else}<span aria-hidden="true"></span>{/if}{/each}
    </div>
  </div>
  <div class="date-picker-footer">
    <button class="button primary full" onclick={apply}
      ><Icon name="calendar" />Tampilkan agenda</button
    >
  </div>
</section>
