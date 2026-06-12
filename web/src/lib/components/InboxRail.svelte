<script lang="ts">
  import Icon from "./Icon.svelte";
  import type { Inbox } from '../types'
  import { inboxes } from '../stores/inboxes.svelte'
  import { messages } from '../stores/messages.svelte'
  import InboxModal from './InboxModal.svelte'

  let modalOpen = $state(false)
  let editing = $state<Inbox | null>(null)

  function openNew() {
    editing = null
    modalOpen = true
  }

  function openEdit(inbox: Inbox) {
    editing = inbox
    modalOpen = true
  }
</script>

<aside class="flex w-80 shrink-0 flex-col overflow-y-auto border-r-2 border-neutral-900 bg-neutral-950">
  <div class="flex items-center border-b-2 border-neutral-900 px-4 py-3">
    <span class="flex items-center gap-2 text-sm font-bold uppercase tracking-wider text-neutral-200">
      <Icon name="inbox" class="text-xl text-accent glow" /> Inboxes
    </span>
    <button onclick={openNew} aria-label="new inbox" class="ml-auto inline-flex items-center justify-center text-neutral-500 hover:text-accent">
      <Icon name="plus" class="text-xl" />
    </button>
  </div>

  <button
    onclick={() => messages.selectInbox(null)}
    class="flex w-full items-center gap-2 border-l-2 px-4 py-2.5 text-left text-sm hover:bg-neutral-900 {messages.activeInbox === null
      ? 'border-l-accent bg-neutral-900 text-accent glow'
      : 'border-l-transparent text-neutral-300'}"
  >
    <Icon name="grid" class="text-xl" />
    <span class="truncate">All</span>
    <span class="ml-auto">[{messages.totalCount}]</span>
  </button>

  {#each inboxes.items as ib (ib.id)}
    <div
      class="flex items-center border-l-2 hover:bg-neutral-900 {messages.activeInbox?.id === ib.id
        ? 'border-l-accent bg-neutral-900'
        : 'border-l-transparent'}"
    >
      <button
        onclick={() => messages.selectInbox(ib)}
        class="flex flex-1 items-center gap-2 py-2.5 pl-4 text-left text-sm {messages.activeInbox?.id ===
        ib.id
          ? 'text-accent glow'
          : 'text-neutral-300'}"
      >
        <Icon name="folder" class="text-xl" />
        <span class="truncate">{ib.name}</span>
        <span class="ml-auto text-neutral-600">[{messages.countForDomain(ib.domain)}]</span>
      </button>
      <button onclick={() => openEdit(ib)} aria-label="manage inbox" class="inline-flex items-center justify-center px-3 py-2.5 text-neutral-600 hover:text-accent">
        <Icon name="cog" class="text-base" />
      </button>
    </div>
  {/each}
</aside>

{#if modalOpen}
  <InboxModal inbox={editing} onclose={() => (modalOpen = false)} />
{/if}
