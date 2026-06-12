<script lang="ts">
  import Icon from './Icon.svelte'
  import Button from './Button.svelte'

  let {
    title = 'Confirm',
    message,
    confirmLabel = 'delete',
    onconfirm,
    oncancel,
  }: {
    title?: string
    message: string
    confirmLabel?: string
    onconfirm: () => void
    oncancel: () => void
  } = $props()
</script>

<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/70">
  <div class="mx-4 w-full max-w-md border-2 border-neutral-900 bg-neutral-950">
    <div class="flex items-center gap-3 border-b-2 border-neutral-900 px-5 py-3">
      <Icon name="trash" class="text-lg text-red-400" />
      <h2 class="text-base font-bold uppercase tracking-widest text-neutral-100">{title}</h2>
      <button
        onclick={oncancel}
        aria-label="close"
        class="ml-auto inline-flex items-center justify-center px-2 text-lg text-neutral-500 hover:text-accent"
      >
        <Icon name="times" />
      </button>
    </div>

    <div class="px-5 py-5 text-sm text-neutral-400">{message}</div>

    <div class="flex justify-end gap-2 border-t-2 border-neutral-900 px-5 py-3">
      <Button onclick={oncancel}>cancel</Button>
      <Button variant="danger" icon="trash" onclick={onconfirm}>{confirmLabel}</Button>
    </div>
  </div>
</div>
