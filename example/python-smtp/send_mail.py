import smtplib
from email.message import EmailMessage

msg = EmailMessage()
msg["Subject"] = "Mail Example From Python"
msg["From"] = "alice@example.com"
msg["To"] = "bob@example.com"

text = """
Hi,
How are you?
This email was sent via Python's smtplib module.
"""

html = """
<html>
    <body>
        <p>Hi,<br>
        How are you?<br>
        This email was sent via Python's smtplib module.
        </p>
    </body>
</html>
"""

try:
    with smtplib.SMTP("localhost", 1025) as smtp:
        msg.set_content(text)
        msg.add_alternative(html, subtype="html")
        smtp.send_message(msg)
    print("Sent - open http://localhost:8025")
except Exception as e:
    print(f"An error occurred: {e}")
