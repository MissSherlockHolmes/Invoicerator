const sgMail = require('@sendgrid/mail');

exports.handler = async (event, context) => {
    if (event.httpMethod !== 'POST') {
        return { statusCode: 405, body: 'Method Not Allowed' };
    }

    try {
        const { recipientEmail, recipientName, companyEmail, companyName, pdfBase64 } = JSON.parse(event.body);

        if (!recipientEmail || !pdfBase64) {
            return {
                statusCode: 400,
                body: JSON.stringify({ error: 'Missing required fields' })
            };
        }

        sgMail.setApiKey(process.env.SENDGRID_API_KEY);

        const senderEmail = process.env.SENDER_EMAIL || 'invoices@invoicerator.com';

        const msg = {
            to: recipientEmail,
            from: {
                email: senderEmail,
                name: companyName || 'Invoicerator'
            },
            subject: `You have received an invoice from ${companyName || 'Invoicerator'}`,
            text: `You have received an invoice from ${companyName || 'Invoicerator'}. Please find it attached to this email.`,
            html: `
                <div style="font-family: Arial, sans-serif; color: #333;">
                    <p>You have received an invoice from ${companyName || 'Invoicerator'}.</p>
                    <p>The invoice is attached to this email as a PDF document.</p>
                    <p>If you have any questions, please contact us directly.</p>
                </div>
            `,
            attachments: [
                {
                    content: pdfBase64,
                    filename: 'invoice.pdf',
                    type: 'application/pdf',
                    disposition: 'attachment'
                }
            ]
        };

        if (companyEmail) {
            msg.cc = companyEmail;
        }

        await sgMail.send(msg);

        return {
            statusCode: 200,
            body: JSON.stringify({ message: 'Invoice sent successfully' })
        };
    } catch (error) {
        console.error('Error sending email:', error);
        return {
            statusCode: 500,
            body: JSON.stringify({ error: 'Failed to send email' })
        };
    }
};
