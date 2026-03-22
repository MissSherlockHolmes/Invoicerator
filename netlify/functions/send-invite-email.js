const sgMail = require('@sendgrid/mail');

exports.handler = async (event, context) => {
    if (event.httpMethod !== 'POST') {
        return { statusCode: 405, body: 'Method Not Allowed' };
    }

    try {
        const { email, token, companyName, origin } = JSON.parse(event.body);

        if (!email || !token) {
            return {
                statusCode: 400,
                body: JSON.stringify({ error: 'Missing required fields' })
            };
        }

        sgMail.setApiKey(process.env.SENDGRID_API_KEY);
        const senderEmail = process.env.SENDER_EMAIL || 'invoices@invoicerator.com';

        const inviteLink = `${origin}/accept_invite.html?token=${token}`;

        const msg = {
            to: email,
            from: {
                email: senderEmail,
                name: 'Invoicerator'
            },
            subject: `You have been invited to join ${companyName} on Invoicerator`,
            text: `You have been invited to join ${companyName} on Invoicerator. Please click the following link to accept the invite: ${inviteLink}`,
            html: `
                <html>
                <body>
                <div style="font-family: Arial, sans-serif; color: #333; max-width: 600px; margin: 0 auto;">
                    <h2>Team Invitation</h2>
                    <p>You have been invited to join <strong>${companyName}</strong> on Invoicerator.</p>
                    <p>By joining, you will be able to create and send invoices on behalf of the company.</p>
                    <p style="margin: 30px 0;">
                        <a href="${inviteLink}" style="background-color: #00bfff; color: white; padding: 12px 24px; text-decoration: none; border-radius: 5px; font-weight: bold;">Accept Invitation</a>
                    </p>
                    <p style="font-size: 0.9em; color: #666;">If the button doesn't work, copy and paste this link into your browser:<br>${inviteLink}</p>
                    <p style="font-size: 0.9em; color: #666;"><em>Note: You must sign up or log in using this exact email address (${email}) to accept the invite.</em></p>
                </div>
                </body>
                </html>
            `
        };

        await sgMail.send(msg);

        return {
            statusCode: 200,
            body: JSON.stringify({ message: 'Invite sent successfully' })
        };
    } catch (error) {
        console.error('Error sending invite email:', error);
        return {
            statusCode: 500,
            body: JSON.stringify({ error: 'Failed to send invite email' })
        };
    }
};
