interface Pending {
  title: string
  message: string
  confirmLabel: string
  onconfirm: () => void
}

class ConfirmStore {
  pending = $state<Pending | null>(null)

  ask = (p: Pending) => {
    this.pending = p
  }

  accept = () => {
    const p = this.pending
    this.pending = null
    p?.onconfirm()
  }

  cancel = () => {
    this.pending = null
  }
}

export const confirm = new ConfirmStore()
