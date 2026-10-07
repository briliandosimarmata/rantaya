<script lang="ts">
  import { getContext } from 'svelte';
  import { goto } from '$app/navigation';
  import type { AppContext } from '#lib/types.js';
  import { mutate, upload } from '#lib/api.js';
  import { returnPath } from '#lib/navigation.js';
  import Avatar from './Avatar.svelte';
  let { onboarding = false, next = '', notifications = false } = $props<{ onboarding?: boolean; next?: string; notifications?: boolean }>();
  const ctx = getContext<AppContext>('app');
  let name = $state(ctx.user?.name || ''),
    city = $state(ctx.user?.city || 'Karawang'),
    bio = $state(ctx.user?.bio || ''),
    avatar = $state(ctx.user?.avatar_url || ''),
    interests = $state<string[]>(ctx.user?.interests || []),
    events = $state(ctx.user?.preferences?.events !== false),
    replies = $state(ctx.user?.preferences?.replies !== false),
    reminders = $state(ctx.user?.preferences?.reminders !== false),
    digest = $state(ctx.user?.preferences?.digest === true),
    busy = $state(false),
    errorMessage = $state('');
  async function save(e?: SubmitEvent, skip = false) {
    e?.preventDefault();
    busy = true;
    try {
      await mutate(
        '/profile',
        {
          name: skip ? ctx.user?.name : name,
          city: skip ? ctx.user?.city : city,
          bio: skip ? ctx.user?.bio : bio,
          avatar_url: avatar,
          interests,
          preferences: { events, replies, reminders, digest },
          onboarding_done: true,
        },
        'PATCH',
      );
      ctx.toast(skip ? 'Profil bisa dilengkapi kapan saja.' : 'Profil tersimpan.');
      if (onboarding)
        await goto(returnPath(next, ctx.user?.role === 'organizer' ? '/kelola' : '/'), { invalidateAll: true });
      else await ctx.refresh();
    } catch (err: any) {
      errorMessage = err.message;
    } finally {
      busy = false;
    }
  }
  async function photo(e: Event) {
    const file = (e.currentTarget as HTMLInputElement).files?.[0];
    if (!file) return;
    busy = true;
    try {
      avatar = (await upload(file)).url;
    } catch (err: any) {
      errorMessage = err.message;
    } finally {
      busy = false;
    }
  }
</script>

<form class="form" onsubmit={(e) => save(e)}>
  {#if !notifications}
  <div class="row">
    <Avatar {name} src={avatar} size="large" />
    <div>
      <label for="avatar">Foto profil</label><input
        id="avatar"
        type="file"
        accept="image/jpeg,image/png,image/webp"
        onchange={photo}
      />
    </div>
  </div>
  <div>
    <label for="name">Nama</label><input
      id="name"
      bind:value={name}
      minlength="2"
      maxlength="80"
      required
    />
  </div>
  <div>
    <label for="city">Kota</label><input
      id="city"
      bind:value={city}
      minlength="2"
      maxlength="60"
      required
    />
  </div>
  <div>
    <label for="bio">Bio singkat</label><textarea id="bio" bind:value={bio} maxlength="500"
    ></textarea>
  </div>
  <fieldset style="border:0;padding:0">
    <legend>Minat</legend>
    <div class="row wrap">
      {#each ['Teater', 'Musik', 'Komedi', 'Seni lainnya'] as value}<label class="inline-label"
          ><input type="checkbox" bind:group={interests} {value} />{value}</label
        >{/each}
    </div>
  </fieldset>
  {/if}
  {#if !onboarding && notifications}<section>
      <div class="check-row"><div><b>Event baru</b><div class="meta">Dari penyelenggara yang kamu ikuti.</div></div><button type="button" role="switch" aria-checked={events} aria-label="Event baru" class="switch" class:on={events} onclick={() => events = !events}></button></div>
      <div class="check-row"><div><b>Balasan percakapan</b><div class="meta">Ketika seseorang membalas komentarmu.</div></div><button type="button" role="switch" aria-checked={replies} aria-label="Balasan percakapan" class="switch" class:on={replies} onclick={() => replies = !replies}></button></div>
      <div class="check-row"><div><b>Pengingat acara</b><div class="meta">Untuk agenda yang kamu simpan.</div></div><button type="button" role="switch" aria-checked={reminders} aria-label="Pengingat acara" class="switch" class:on={reminders} onclick={() => reminders = !reminders}></button></div>
      <div class="check-row"><div><b>Rangkuman email</b><div class="meta">Rangkuman saat ada kabar baru.</div></div><button type="button" role="switch" aria-checked={digest} aria-label="Rangkuman email" class="switch" class:on={digest} onclick={() => digest = !digest}></button></div>
      <p class="meta" style="margin-top:18px">Pilihan email disimpan; pengiriman email belum tersedia. Status pembayaran dan tiket tetap dikabarkan.</p>
    </section>{/if}{#if errorMessage}<p class="note error-note" role="alert">
      {errorMessage}
    </p>{/if}<button class="button primary" disabled={busy}
    >{busy ? 'Menyimpan…' : notifications ? 'Simpan pengaturan' : onboarding ? 'Simpan & lanjut' : 'Simpan profil'}</button
  >{#if onboarding}<button
      class="action"
      type="button"
      onclick={() => save(undefined, true)}
      disabled={busy}>Lewati, lengkapi nanti</button
    >{/if}
</form>
