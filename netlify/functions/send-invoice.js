const sgMail = require('@sendgrid/mail');

exports.handler = async (event, context) => {
    if (event.httpMethod !== 'POST') {
        return { statusCode: 405, body: 'Method Not Allowed' };
    }

    try {
        const { recipientEmail, recipientName, companyEmail, companyName, pdfBase64, ccEmails, bccEmails, filename } = JSON.parse(event.body);

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
            text: `You have received an invoice from ${companyName || 'Invoicerator'}. Please find it attached to this email as a PDF document. If you have any questions, please contact us directly.`,
            html: `
                <html>
                <body>
                <div style="font-family: Arial, sans-serif; color: #333;">
                    <p>You have received an invoice from ${companyName || 'Invoicerator'}.</p>
                    <p>The invoice is attached to this email as a PDF document.</p>
                    <p>If you have any questions, please contact us directly.</p>
                </div>
                </body>
                </html>
            `,
            attachments: [
                {
                    content: pdfBase64,
                    filename: filename || 'invoice.pdf',
                    type: 'application/pdf',
                    disposition: 'attachment'
                }
            ]
        };

        let ccArray = [];
        if (companyEmail) {
            ccArray.push(companyEmail);
        }
        if (ccEmails) {
            const parsedCcs = ccEmails.split(',').map(e => e.trim()).filter(e => e);
            ccArray = ccArray.concat(parsedCcs);
        }
        
        // Remove duplicates and ensure we don't CC the recipient or sender
        ccArray = [...new Set(ccArray)].filter(e => e !== recipientEmail && e !== senderEmail);

        if (ccArray.length > 0) {
            msg.cc = ccArray;
        }

        let bccArray = [];
        if (bccEmails) {
            const parsedBccs = bccEmails.split(',').map(e => e.trim()).filter(e => e);
            bccArray = bccArray.concat(parsedBccs);
        }

        // Remove duplicates and ensure we don't BCC the recipient, sender, or someone already in CC
        bccArray = [...new Set(bccArray)].filter(e => e !== recipientEmail && e !== senderEmail && !ccArray.includes(e));

        if (bccArray.length > 0) {
            msg.bcc = bccArray;
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
