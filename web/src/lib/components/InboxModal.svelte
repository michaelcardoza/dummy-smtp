<script lang="ts">
  import Icon from "./Icon.svelte";
  import type { Inbox } from '../types'
  import { inboxes } from '../stores/inboxes.svelte'
  import { messages } from '../stores/messages.svelte'

  let { inbox, onclose }: { inbox: Inbox | null; onclose: () => void } = $props()

  // svelte-ignore state_referenced_locally
  let name = $state(inbox?.name ?? '')
  // svelte-ignore state_referenced_locally
  let domain = $state(inbox?.domain ?? '')
  // svelte-ignore state_referenced_locally
  let description = $state(inbox?.description ?? '')

  function save() {
    if (!name.trim() || !domain.trim()) return
    const data = { name: name.trim(), domain: domain.trim(), description: description.trim() }
    if (inbox) inboxes.update(inbox.id, data)
    else inboxes.create(data)
    onclose()
  }

  function del() {
    if (!inbox) return
    if (messages.activeInbox?.id === inbox.id) messages.selectInbox(null)
    inboxes.remove(inbox.id)
    onclose()
  }
</script>

<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/70">
  <div class="mx-4 w-full max-w-xl border-2 border-neutral-900 bg-neutral-950">
    <div class="flex items-center gap-3 border-b-2 border-neutral-900 px-5 py-3">
      <Icon name="folder" class="glow text-lg text-accent" />
      <h2 class="text-base font-bold uppercase tracking-widest text-neutral-100">
        <span class="text-accent glow">&gt;</span> {inbox ? 'Edit inbox' : 'New inbox'}
      </h2>
      <button onclick={onclose} aria-label="close" class="ml-auto inline-flex items-center justify-center px-2 text-lg text-neutral-500 hover:text-accent">
        <Icon name="times" />
      </button>
    </div>

    <div class="space-y-4 px-5 py-4">
      <label class="block">
        <span class="mb-1 block text-sm text-neutral-400">Name</span>
        <input
          bind:value={name}
          placeholder="Project A"
          class="w-full border-2 border-neutral-900 bg-neutral-950 px-3 py-2 text-sm text-neutral-200 placeholder:text-neutral-600 focus:border-accent focus:outline-none"
        />
      </label>

      <label class="block">
        <span class="mb-1 block text-sm text-neutral-400">Domain</span>
        <input
          bind:value={domain}
          placeholder="projecta.test"
          class="w-full border-2 border-neutral-900 bg-neutral-950 px-3 py-2 text-sm text-neutral-200 placeholder:text-neutral-600 focus:border-accent focus:outline-none"
        />
      </label>

      <label class="block">
        <span class="mb-1 block text-sm text-neutral-400">Description</span>
        <input
          bind:value={description}
          placeholder="optional"
          class="w-full border-2 border-neutral-900 bg-neutral-950 px-3 py-2 text-sm text-neutral-200 placeholder:text-neutral-600 focus:border-accent focus:outline-none"
        />
      </label>

      {#if domain.trim()}
        <div class="border-2 border-neutral-900 bg-neutral-950 p-4 text-sm">
          <div class="mb-2 flex items-center gap-2 text-neutral-400">
            <Icon name="info-circle" class="glow text-accent" /> How mail reaches this inbox
          </div>
          <p class="mb-2 text-neutral-500">
            Any message whose recipient uses
            <span class="text-accent glow">@{domain.trim()}</span> lands here.
          </p>
          <pre class="overflow-x-auto border-2 border-neutral-900 bg-neutral-950 p-3 text-neutral-300">To: user@{domain.trim()}
SMTP: localhost:1025</pre>
        </div>
      {/if}
    </div>

    <div class="flex items-center gap-2 border-t-2 border-neutral-900 px-5 py-3">
      {#if inbox}
        <button
          onclick={del}
          class="flex items-center gap-2 border-2 border-red-900/60 px-3 py-1.5 text-sm text-red-400 hover:bg-red-500 hover:text-white"
        >
          <Icon name="trash" /> delete
        </button>
      {/if}
      <button
        onclick={save}
        class="ml-auto flex items-center gap-2 border-2 border-accent/40 px-4 py-1.5 text-sm text-accent hover:bg-accent hover:text-black"
      >
        <Icon name="check-circle" /> save
      </button>
    </div>
  </div>
</div>
