<script lang="ts">
  import Icon from "./lib/components/Icon.svelte";
  import { onMount } from 'svelte'
  import { messages } from './lib/stores/messages.svelte'
  import { confirm } from './lib/stores/confirm.svelte'
  import { SMTP_HOST, SMTP_PORT } from './lib/config'
  import Header from './lib/components/Header.svelte'
  import MessageList from './lib/components/MessageList.svelte'
  import MessageDetail from './lib/components/MessageDetail.svelte'
  import Toaster from './lib/components/Toaster.svelte'
  import ConfirmModal from './lib/components/ConfirmModal.svelte'

  onMount(() => {
    messages.load()
    messages.connectLive()
    return () => messages.disconnectLive()
  })
</script>

<div class="flex h-screen flex-col bg-neutral-950 font-mono text-base text-neutral-200 antialiased">
  <Header />
  <main class="flex flex-1 overflow-hidden">
    <MessageList />
    <section class="flex-1 overflow-y-auto bg-neutral-950">
      {#if messages.selected}
        <MessageDetail message={messages.selected} />
      {:else}
        <div class="flex h-full flex-col items-center justify-center gap-3 text-neutral-500">
          <Icon name="envelope" class="text-4xl text-neutral-700" />
          <div>Select a message to view it.<span class="blink text-accent">_</span></div>
          <div class="text-sm text-neutral-600">or send one to <span class="text-accent glow">{SMTP_HOST}:{SMTP_PORT}</span></div>
        </div>
      {/if}
    </section>
  </main>
</div>

<Toaster />

{#if confirm.pending}
  <ConfirmModal
    title={confirm.pending.title}
    message={confirm.pending.message}
    confirmLabel={confirm.pending.confirmLabel}
    onconfirm={confirm.accept}
    oncancel={confirm.cancel}
  />
{/if}
