import smtplib
from email.message import EmailMessage

msg = EmailMessage()
msg["Subject"] = "Order #10492 confirmed - Acme Store"
msg["From"] = "Acme Store <orders@acme.example.com>"
msg["To"] = "Jane Doe <jane.doe@example.com>"
msg["Cc"] = "billing@example.com, John Doe <john.doe@example.com>"
msg["Reply-To"] = "support@acme.example.com"

text = """
Hi Jane,

Thanks for your order! Here's a summary.

Order:      #10492
Date:       2026-06-14
Items:
  - Mechanical Keyboard (x1)   $89.00
  - USB-C Cable 2m (x2)        $18.00
Subtotal:   $107.00
Shipping:   $5.00
Total:      $112.00

Shipping to:
  Jane Doe
  742 Evergreen Terrace
  Springfield

The invoice is attached as a CSV.

- Acme Store
"""

html = """
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="color-scheme" content="light dark">
    <meta name="supported-color-schemes" content="light dark">
    <title>Order #10492 confirmed</title>
    <style>
      :root { color-scheme: light dark; supported-color-schemes: light dark; }
      @media (prefers-color-scheme: dark) {
        .email-bg   { background: #09090b !important; }
        .card       { background: #18181b !important; border-color: #27272a !important; }
        .panel      { background: #1c1c1f !important; border-color: #27272a !important; }
        .footer     { background: #131316 !important; border-color: #27272a !important; }
        .heading    { color: #fafafa !important; }
        .body-text  { color: #d4d4d8 !important; }
        .muted      { color: #a1a1aa !important; }
        .faint      { color: #71717a !important; }
        .divider    { border-color: #27272a !important; }
        .total-row  { border-color: #3f3f46 !important; }
      }
    </style>
  </head>
  <body class="email-bg" style="font-family: -apple-system, 'Segoe UI', Arial, sans-serif; color: #222; margin: 0; padding: 32px 16px; background: #f4f4f5;">
    <table class="card" width="100%" cellpadding="0" cellspacing="0" style="max-width: 600px; margin: 0 auto; background: #fff; border-radius: 12px; overflow: hidden; border: 1px solid #e4e4e7;">
      <tr>
        <td style="background: #18181b; padding: 24px 32px;">
          <table width="100%" cellpadding="0" cellspacing="0">
            <tr>
              <td style="color: #fff; font-size: 20px; font-weight: bold; letter-spacing: 0.5px;">ACME STORE</td>
              <td align="right" style="color: #a1a1aa; font-size: 13px;">Order #10492</td>
            </tr>
          </table>
        </td>
      </tr>

      <tr>
        <td style="padding: 0;">
          <img src="https://placehold.co/600x200/6366f1/ffffff/png?text=Thanks+for+your+order" width="600" alt="Thanks for your order" style="display: block; width: 100%; height: auto;">
        </td>
      </tr>

      <tr>
        <td style="padding: 32px 32px 8px;">
          <h1 class="heading" style="margin: 0 0 8px; font-size: 24px; color: #18181b;">Your order is confirmed 🎉</h1>
          <p class="body-text" style="margin: 0; color: #555; font-size: 15px; line-height: 1.6;">
            Hi Jane, thanks for shopping with us! We received your order on
            <strong>June 14, 2026</strong> and we're getting it ready to ship.
            You'll get a tracking link as soon as it leaves our warehouse.
          </p>
        </td>
      </tr>

      <tr>
        <td style="padding: 24px 32px 8px;">
          <table class="panel" width="100%" cellpadding="0" cellspacing="0" style="border: 1px solid #e4e4e7; border-radius: 8px; overflow: hidden;">
            <tr class="panel" style="background: #fafafa;">
              <td class="faint" style="padding: 12px 16px; font-size: 12px; text-transform: uppercase; letter-spacing: 0.5px; color: #71717a;">Order summary</td>
            </tr>

            <tr>
              <td class="divider" style="padding: 16px; border-top: 1px solid #f0f0f0;">
                <table cellpadding="0" cellspacing="0">
                  <tr>
                    <td style="padding-right: 14px;">
                      <img src="https://placehold.co/56x56/e4e4e7/71717a/png?text=KB" width="56" height="56" alt="Mechanical Keyboard" style="display: block; border-radius: 6px;">
                    </td>
                    <td>
                      <div class="heading" style="font-size: 14px; font-weight: 600; color: #18181b;">Mechanical Keyboard</div>
                      <div class="faint" style="font-size: 13px; color: #71717a; margin-top: 2px;">Brown switches · Qty 1</div>
                    </td>
                  </tr>
                </table>
              </td>
              <td class="heading divider" align="right" style="padding: 16px; border-top: 1px solid #f0f0f0; font-size: 14px; font-weight: 600; color: #18181b; vertical-align: top;">$89.00</td>
            </tr>

            <tr>
              <td class="divider" style="padding: 16px; border-top: 1px solid #f0f0f0;">
                <table cellpadding="0" cellspacing="0">
                  <tr>
                    <td style="padding-right: 14px;">
                      <img src="https://placehold.co/56x56/e4e4e7/71717a/png?text=USB" width="56" height="56" alt="USB-C Cable" style="display: block; border-radius: 6px;">
                    </td>
                    <td>
                      <div class="heading" style="font-size: 14px; font-weight: 600; color: #18181b;">USB-C Cable 2m</div>
                      <div class="faint" style="font-size: 13px; color: #71717a; margin-top: 2px;">Braided · Qty 2</div>
                    </td>
                  </tr>
                </table>
              </td>
              <td class="heading divider" align="right" style="padding: 16px; border-top: 1px solid #f0f0f0; font-size: 14px; font-weight: 600; color: #18181b; vertical-align: top;">$18.00</td>
            </tr>

            <tr>
              <td class="faint divider" style="padding: 10px 16px; border-top: 1px solid #f0f0f0; font-size: 14px; color: #71717a;">Subtotal</td>
              <td class="faint divider" align="right" style="padding: 10px 16px; border-top: 1px solid #f0f0f0; font-size: 14px; color: #71717a;">$107.00</td>
            </tr>
            <tr>
              <td class="faint" style="padding: 10px 16px; font-size: 14px; color: #71717a;">Shipping</td>
              <td class="faint" align="right" style="padding: 10px 16px; font-size: 14px; color: #71717a;">$5.00</td>
            </tr>
            <tr>
              <td class="heading total-row" style="padding: 14px 16px; border-top: 2px solid #e4e4e7; font-size: 16px; font-weight: bold; color: #18181b;">Total</td>
              <td class="heading total-row" align="right" style="padding: 14px 16px; border-top: 2px solid #e4e4e7; font-size: 16px; font-weight: bold; color: #18181b;">$112.00</td>
            </tr>
          </table>
        </td>
      </tr>

      <tr>
        <td style="padding: 24px 32px;">
          <table width="100%" cellpadding="0" cellspacing="0">
            <tr>
              <td width="50%" style="vertical-align: top; padding-right: 8px;">
                <div class="faint" style="font-size: 12px; text-transform: uppercase; letter-spacing: 0.5px; color: #71717a; margin-bottom: 6px;">Shipping address</div>
                <div class="body-text" style="font-size: 14px; color: #333; line-height: 1.6;">
                  Jane Doe<br>742 Evergreen Terrace<br>Springfield, IL 62704<br>United States
                </div>
              </td>
              <td width="50%" style="vertical-align: top; padding-left: 8px;">
                <div class="faint" style="font-size: 12px; text-transform: uppercase; letter-spacing: 0.5px; color: #71717a; margin-bottom: 6px;">Delivery estimate</div>
                <div class="body-text" style="font-size: 14px; color: #333; line-height: 1.6;">
                  Standard shipping<br><strong>June 18 – June 20</strong>
                </div>
              </td>
            </tr>
          </table>
        </td>
      </tr>

      <tr>
        <td align="center" style="padding: 8px 32px 32px;">
          <a href="https://acme.example.com/orders/10492" style="background: #6366f1; color: #fff; text-decoration: none; padding: 14px 32px; border-radius: 8px; display: inline-block; font-size: 15px; font-weight: 600;">Track your order</a>
          <p class="faint" style="margin: 16px 0 0; font-size: 13px; color: #71717a;">
            or <a href="https://acme.example.com/orders/10492/invoice" style="color: #6366f1;">view your invoice online</a>
          </p>
        </td>
      </tr>

      <tr>
        <td class="footer" style="background: #fafafa; border-top: 1px solid #e4e4e7; padding: 24px 32px;">
          <p class="heading" style="margin: 0 0 12px; font-size: 14px; font-weight: 600; color: #18181b;">Need help?</p>
          <p class="faint" style="margin: 0 0 16px; font-size: 13px; color: #71717a; line-height: 1.6;">
            Reply to this email or reach us at
            <a href="mailto:support@acme.example.com" style="color: #6366f1;">support@acme.example.com</a>.
            The full invoice is attached as a CSV.
          </p>
          <table cellpadding="0" cellspacing="0">
            <tr>
              <td style="padding-right: 12px;"><a href="#" style="color: #71717a; text-decoration: none; font-size: 13px;">Twitter</a></td>
              <td style="padding-right: 12px;"><a href="#" style="color: #71717a; text-decoration: none; font-size: 13px;">Instagram</a></td>
              <td><a href="#" style="color: #71717a; text-decoration: none; font-size: 13px;">Help center</a></td>
            </tr>
          </table>
        </td>
      </tr>

      <tr>
        <td align="center" style="padding: 20px 32px; color: #a1a1aa; font-size: 12px; line-height: 1.6;">
          Acme Store Inc. · 500 Market St, San Francisco, CA<br>
          You're receiving this because you placed an order.
          <a href="#" style="color: #a1a1aa;">Unsubscribe</a>
        </td>
      </tr>
    </table>
  </body>
</html>
"""

invoice_csv = (
    "item,qty,unit_price,total\n"
    "Mechanical Keyboard,1,89.00,89.00\n"
    "USB-C Cable 2m,2,9.00,18.00\n"
    ",,subtotal,107.00\n"
    ",,shipping,5.00\n"
    ",,total,112.00\n"
)

try:
    msg.set_content(text)
    msg.add_alternative(html, subtype="html")
    msg.add_attachment(
        invoice_csv.encode("utf-8"),
        maintype="text",
        subtype="csv",
        filename="invoice-10492.csv",
    )

    with smtplib.SMTP("localhost", 1025) as smtp:
        smtp.send_message(msg)
    print("Sent - open http://localhost:8025")
except Exception as e:
    print(f"An error occurred: {e}")
