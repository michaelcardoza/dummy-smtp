<script lang="ts">
  import type { Snippet } from 'svelte'
  import Icon, { type IconName } from './Icon.svelte'

  type Variant = 'neutral' | 'accent' | 'danger'
  type Size = 'sm' | 'md'

  let {
    variant = 'neutral',
    size = 'sm',
    icon,
    href,
    target,
    title,
    onclick,
    children,
  }: {
    variant?: Variant
    size?: Size
    icon?: IconName
    href?: string
    target?: string
    title?: string
    onclick?: () => void
    children?: Snippet
  } = $props()

  const VARIANTS: Record<Variant, string> = {
    neutral: 'border-neutral-900 text-neutral-400 hover:border-neutral-600 hover:text-neutral-100',
    accent: 'border-accent/40 text-accent hover:bg-accent hover:text-black',
    danger: 'border-red-900/60 text-red-400 hover:bg-red-500 hover:text-white',
  }
  const SIZES: Record<Size, string> = {
    sm: 'px-3 py-1.5',
    md: 'h-10 px-3',
  }

  let className = $derived(
    `flex items-center gap-2 border-2 text-sm ${SIZES[size]} ${VARIANTS[variant]}`,
  )
</script>

{#if href}
  <a
    {href}
    {target}
    {title}
    rel={target === '_blank' ? 'noreferrer' : undefined}
    class={className}
  >
    {#if icon}<Icon name={icon} />{/if}
    {@render children?.()}
  </a>
{:else}
  <button {onclick} {title} class={className}>
    {#if icon}<Icon name={icon} />{/if}
    {@render children?.()}
  </button>
{/if}
