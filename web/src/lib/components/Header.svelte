<script lang="ts">
  import Icon from './Icon.svelte'
  import Button from './Button.svelte'
  import { messages } from '../stores/messages.svelte'
  import { confirm } from '../stores/confirm.svelte'
  import { REPO_URL } from '../config'
  import ConnectHelp from './ConnectHelp.svelte'

  let showConnect = $state(false)

  let liveState = $derived(
    !messages.liveEnabled ? 'paused' : messages.isConnected ? 'live' : 'offline',
  )

  let filtering = $derived(messages.searchQuery.trim() !== '')

  function confirmClear() {
    confirm.ask({
      title: 'Clear all',
      message: "Delete all messages? This can't be undone.",
      confirmLabel: 'clear all',
      onconfirm: messages.clear,
    })
  }
</script>

<header class="flex items-center gap-4 border-b-2 border-neutral-900 bg-neutral-950 px-6 py-4">
  <Icon name="envelope" class="glow text-2xl text-accent" />
  <h1 class="text-lg font-bold uppercase tracking-widest text-neutral-100">
    <span class="text-accent glow">&gt;</span> Dummy_SMTP
  </h1>

  <span class="flex items-center gap-2 text-sm font-bold tracking-wider {liveState === 'live'
    ? 'text-accent glow'
    : liveState === 'offline'
      ? 'text-red-400'
      : 'text-neutral-500'}">
    <span class="h-3 w-3 {liveState === 'live'
      ? 'bg-accent blink'
      : liveState === 'offline'
        ? 'bg-red-500'
        : 'bg-neutral-600'}"></span>
    {liveState === 'live' ? 'LIVE' : liveState === 'offline' ? 'OFFLINE' : 'PAUSED'}
  </span>
  <button
    onclick={messages.toggleLive}
    title={messages.liveEnabled ? 'Pause live updates' : 'Resume live updates'}
    aria-label="toggle live updates"
    class="inline-flex items-center justify-center {messages.liveEnabled
      ? 'text-accent hover:text-neutral-200'
      : 'text-neutral-500 hover:text-accent'}"
  >
    <Icon name="power" class="text-lg" />
  </button>

  <!--<div class="relative ml-4 w-80">
    <Icon name="search" class="absolute left-2.5 top-1/2 -translate-y-1/2 text-neutral-600" />
    <input
      type="text"
      placeholder="search…"
      bind:value={messages.searchQuery}
      class="h-10 w-full border-2 border-neutral-900 bg-neutral-950 pl-8 pr-3 text-sm text-neutral-200 placeholder:text-neutral-600 focus:border-accent focus:outline-none"
    />
  </div>-->

  <div class="ml-auto flex h-10 items-center gap-2 bg-accent/10 px-3 text-sm">
    <span class="font-bold tabular-nums text-accent"><span class="text-lg">[</span>{messages.filteredCount}{#if filtering}<span class="text-accent/50">/{messages.totalCount}</span>{/if}<span class="text-lg">]</span></span>
    <span class="text-xs uppercase tracking-wider text-accent/70">msgs</span>
  </div>

  <Button variant="accent" size="md" icon="code" onclick={() => (showConnect = true)}>connect</Button>
  <button
    onclick={messages.sendTest}
    class="flex items-center gap-2 border-2 border-neutral-900 h-10 px-3 text-sm text-neutral-400 hover:border-accent hover:text-accent"
  >
    <Icon name="send" /> send test
  </button>
  <Button size="md" icon="refresh" onclick={messages.load}>refresh</Button>
  <Button variant="danger" size="md" icon="trash" onclick={confirmClear}>clear</Button>
  <a
    href={REPO_URL}
    target="_blank"
    rel="noreferrer"
    title="GitHub"
    aria-label="GitHub repository"
    class="inline-flex h-10 items-center justify-center border-2 border-neutral-900 px-2.5 text-neutral-400 hover:border-neutral-500 hover:text-white"
  >
    <Icon name="github" class="text-lg" />
  </a>
  <!--<button
    class="inline-flex h-10 items-center justify-center border-2 border-neutral-900 px-2.5 text-sm text-neutral-500 hover:border-neutral-600 hover:text-neutral-200"
    title="settings"
  >
    <Icon name="cog" />
  </button>-->
</header>

{#if showConnect}
  <ConnectHelp onclose={() => (showConnect = false)} />
{/if}
