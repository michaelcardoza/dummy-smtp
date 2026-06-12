<script lang="ts">
  import Icon from './Icon.svelte'
  import { copy } from '../utils/clipboard'
  import { highlight } from '../utils/highlight'
  import { SMTP_HOST, SMTP_PORT } from '../config'

  let { onclose }: { onclose: () => void } = $props()

  let copied = $state(false)

  async function copySnippet() {
    await copy(snippets[lang])
    copied = true
    setTimeout(() => (copied = false), 1500)
  }

  type Lang = 'python' | 'go' | 'node' | 'php' | 'ruby' | 'curl'
  let lang = $state<Lang>('python')

  const langs: [Lang, string][] = [
    ['python', 'Python'],
    ['go', 'Go'],
    ['node', 'Node'],
    ['php', 'PHP'],
    ['ruby', 'Ruby'],
    ['curl', 'curl'],
  ]

  const grammarByLang: Record<Lang, string> = {
    python: 'python',
    go: 'go',
    node: 'javascript',
    php: 'php',
    ruby: 'ruby',
    curl: 'bash',
  }

  const snippets: Record<Lang, string> = {
    python: `import smtplib
from email.message import EmailMessage

msg = EmailMessage()
msg["Subject"] = "Mail Example From Python"
msg["From"] = "alice@example.com"
msg["To"] = "bob@example.com"

text = "Hi, how are you? This email was sent via Python's smtplib."
html = "<p>Hi, how are you? This email was sent via Python's smtplib.</p>"

with smtplib.SMTP("localhost", 1025) as smtp:
    msg.set_content(text)
    msg.add_alternative(html, subtype="html")
    smtp.send_message(msg)`,
    go: `package main

import gomail "gopkg.in/mail.v2"

func main() {
	m := gomail.NewMessage()
	m.SetHeader("From", "alice@example.com")
	m.SetHeader("To", "bob@example.com")
	m.SetHeader("Subject", "Mail Example From Go")
	m.SetBody("text/plain", "Hi, how are you? This email was sent via Go.")
	m.AddAlternative("text/html", "<p>Hi, how are you? This email was sent via Go.</p>")

	d := gomail.NewDialer("localhost", 1025, "", "")
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}`,
    node: `import nodemailer from "nodemailer";

const transport = nodemailer.createTransport({
  host: "localhost",
  port: 1025,
  secure: false,
});

await transport.sendMail({
  from: "alice@example.com",
  to: "bob@example.com",
  subject: "Mail Example From JavaScript",
  text: "Hi, how are you? This email was sent via Nodemailer.",
  html: "<p>Hi, how are you? This email was sent via Nodemailer.</p>",
});`,
    php: `<?php
use PHPMailer\\PHPMailer\\PHPMailer;
require 'vendor/autoload.php';

$mail = new PHPMailer();
$mail->isSMTP();
$mail->Host = 'localhost';
$mail->Port = 1025;
$mail->setFrom('alice@example.com');
$mail->addAddress('bob@example.com');
$mail->Subject = 'Mail Example From PHP';
$mail->Body = 'Hi, how are you? This email was sent via PHPMailer.';
$mail->send();`,
    ruby: `require 'mail'

Mail.defaults do
  delivery_method :smtp, address: 'localhost', port: 1025
end

Mail.deliver do
  from    'alice@example.com'
  to      'bob@example.com'
  subject 'Mail Example From Ruby'
  body    'Hi, how are you? This email was sent via the Mail gem.'
end`,
    curl: `curl smtp://localhost:1025 \\
  --mail-from alice@example.com \\
  --mail-rcpt bob@example.com \\
  -T <(printf 'Subject: Mail Example From curl\\n\\nHi, how are you? This email was sent via curl.\\n')`,
  }

  let highlighted = $derived(highlight(snippets[lang], grammarByLang[lang]))
</script>

<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/70">
  <div class="mx-4 w-full max-w-xl border-2 border-neutral-900 bg-neutral-950">
    <div class="flex items-center gap-3 border-b-2 border-neutral-900 px-5 py-3">
      <Icon name="envelope" class="glow text-lg text-accent" />
      <h2 class="text-base font-bold uppercase tracking-widest text-neutral-100">
        <span class="text-accent glow">&gt;</span> How to connect
      </h2>
      <button onclick={onclose} aria-label="close" class="ml-auto inline-flex items-center justify-center px-2 text-lg text-neutral-500 hover:text-accent">
        <Icon name="times" />
      </button>
    </div>

    <div class="px-5 py-4">
      <div class="mb-5 border-2 border-neutral-900 bg-neutral-950 px-4 py-3 text-sm text-neutral-400">
        host: <span class="text-accent glow">{SMTP_HOST}</span> ·
        port: <span class="text-accent glow">{SMTP_PORT}</span> ·
        no auth · no TLS
      </div>

      <div class="mb-3 flex gap-2 text-sm">
        {#each langs as [id, label]}
          <button
            onclick={() => (lang = id)}
            class="border-2 px-3 py-1.5 {lang === id
              ? 'border-accent text-accent glow'
              : 'border-transparent text-neutral-500 hover:text-neutral-300'}"
          >
            {label}
          </button>
        {/each}
      </div>

      <div class="relative">
        <button
          onclick={copySnippet}
          class="absolute right-2 top-2 flex items-center gap-1.5 border-2 px-2 py-1 text-xs {copied
            ? 'border-accent text-accent'
            : 'border-neutral-900 bg-neutral-950 text-neutral-400 hover:border-accent hover:text-accent'}"
        >
          <Icon name={copied ? 'check-circle' : 'copy'} />
          {copied ? 'copied' : 'copy'}
        </button>
        <pre class="max-h-96 overflow-y-auto border-2 border-neutral-900 bg-neutral-950 p-4 text-sm leading-relaxed text-neutral-300">{@html highlighted}</pre>
      </div>
    </div>
  </div>
</div>
