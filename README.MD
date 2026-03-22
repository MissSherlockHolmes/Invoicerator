**Invoicerator**
================

Invoicerator is a streamlined web application designed for creating and sending custom PDF invoices with ease. The app prioritizes user privacy by ensuring that invoices are deleted immediately after being sent. Additionally, users can upload previously generated invoices to auto-populate fields and edit them effortlessly.

Check it out here: [Invoicerator.com](https://invoicerator.com/)

* * * * *

**Features**
------------

-   **User Authentication**: Secure login and signup functionality with session management using Supabase.
-   **Invoice Creation**: Generate professional invoices with customizable fields and line items.
-   **Invoice Preview**: Preview invoices before finalizing or sending them directly in the browser.
-   **PDF Generation**: Automatically create PDFs of invoices, including company letterhead and branding, completely client-side.
-   **Email Integration**: Send invoices directly via email using SendGrid via Netlify Functions.
-   **Profile Management**: Personalize company details, letterheads, and invoice settings.
-   **Secure Deletion**: Invoices are generated locally and sent via email, never stored on a server.

* * * * *

**Project Overview**
--------------------

| **Category** | **Details** |
| --- | --- |
| **Primary Function** | Create, send, and manage invoices securely and efficiently. |
| **Tech Stack** | HTML/CSS/JS, Supabase (Auth/DB/Storage), Netlify Functions (Node.js) |
| **Deployment** | Hosted on [Netlify](https://www.netlify.com/). |
| **Database** | Supabase PostgreSQL for user profiles and settings. |
| **Target Audience** | Small businesses, freelancers, and professionals needing invoicing tools. |

* * * * *

**Setup and Installation**
--------------------------

### **1. Prerequisites**

-   Node.js and npm
-   Netlify CLI (`npm install -g netlify-cli`)
-   A Supabase account
-   A SendGrid account

### **2. Clone the Repository**
```bash
git clone https://github.com/<your-username>/Invoicerator.git
cd Invoicerator
```

### **3. Install Dependencies**
```bash
npm install
```

### **4. Configure Supabase**
1. Create a new project in Supabase.
2. Run the SQL commands in `supabase_schema.sql` in your Supabase SQL Editor to create the necessary tables and policies.
3. Create a storage bucket named `letterheads` and make it public.
4. Update `public/js/supabase.js` with your Supabase Project URL and Anon Key.

### **5. Configure Environment Variables**
Create a `.env` file in the root directory for Netlify Functions:
```
SENDGRID_API_KEY=<your-sendgrid-api-key>
```

### **6. Run the Application Locally**

```bash
netlify dev
```

The application will be available at `http://localhost:8888`.

* * * * *

**Contributors**
----------------

-   **Amberleigh Wankel**
-   **Maksim Azarskov**
-   **Sanskar Chaudhary**
-   **Reika Yonetsu**

* * * * *

**Contact and Support**
-----------------------

For any inquiries or support, please email us at: support@invoicerator.com.

* * * * *
