import nodemailer from "nodemailer";

const transport = nodemailer.createTransport({
  host: "localhost",
  port: 1025,
  secure: false,
});

const text = `
Hi,
How are you?
This email was sent via Nodemailer.
`;

const html = `
<html>
    <body>
        <p>Hi,<br>
        How are you?<br>
        This email was sent via Nodemailer.
        </p>
    </body>
</html>
`;

try {
  await transport.sendMail({
    from: "alice@example.com",
    to: "bob@example.com",
    subject: "Mail Example From JavaScript",
    text,
    html,
  });
  console.log("Sent - open http://localhost:8025");
} catch (e) {
  console.error("An error occurred:", e.message);
}
