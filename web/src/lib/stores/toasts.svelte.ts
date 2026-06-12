interface Toast {
  id: number
  text: string
  error: boolean
}

class ToastStore {
  items = $state<Toast[]>([])

  push = (text: string, error = false) => {
    const id = Date.now() + Math.random()
    this.items = [...this.items, { id, text, error }]
    setTimeout(() => this.dismiss(id), 4000)
  }

  dismiss = (id: number) => {
    this.items = this.items.filter((t) => t.id !== id)
  }
}

export const toasts = new ToastStore()
