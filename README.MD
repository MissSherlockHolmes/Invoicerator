**Invoicerator**
================

Invoicerator is a streamlined web application designed for creating and sending custom PDF invoices with ease. The app prioritizes user privacy by ensuring that invoices are deleted immediately after being sent. Additionally, users can upload previously generated invoices to auto-populate fields and edit them effortlessly.

* * * * *

**Features**
------------

-   **User Authentication**: Secure login and signup functionality with session management.
-   **Invoice Creation**: Generate professional invoices with customizable fields and line items.
-   **Invoice Preview**: Preview invoices before finalizing or sending them.
-   **PDF Generation**: Automatically create PDFs of invoices, including company letterhead and branding.
-   **Email Integration**: Send invoices directly via email using SendGrid.
-   **Profile Management**: Personalize company details, letterheads, and invoice settings.
-   **Secure Deletion**: Invoices are deleted immediately after being sent.

* * * * *

**Project Overview**
--------------------

| **Category** | **Details** |
| --- | --- |
| **Primary Function** | Create, send, and manage invoices securely and efficiently. |
| **Tech Stack** | Go (Gin Framework), SQLite, HTML/CSS, Bootstrap |
| **Deployment** | Hosted on [AWS EC2](https://aws.amazon.com/ec2/) with NGINX and SSL (Let's Encrypt). |
| **Database** | SQLite for simplicity and reliability. |
| **Target Audience** | Small businesses, freelancers, and professionals needing invoicing tools. |

* * * * *

**Setup and Installation**
--------------------------

### **1\. Prerequisites**

-   Go v1.20 or higher
-   Git
-   SQLite
-   An AWS EC2 instance (or local server) for deployment

### **2\. Clone the Repository**
Copy code
    
    git clone https://github.com/<your-username>/Invoicerator.git
    cd Invoicerator
    

### **3\. Configure Environment Variables**

Create the following `.env` files:

-   **Local Development (`.env.local`)**:
    ```
    ENV=local
    DATABASE_URL=invoicerator.db
    SENDGRID_API_KEY=<your-sendgrid-api-key>
    ```

-   **Production (`.env.production`)**:

    ```
    ENV=production
    DATABASE_URL=/app/invoicerator.db
    SENDGRID_API_KEY=<your-sendgrid-api-key>
    ```

### **4\. Run the Application Locally**

`go run main.go`

The application will be available at `http://localhost:8080`.

* * * * *


**Contributors**
----------------

-   **Amberleigh Wankel**
-   **Maksim Azarskov**
-   **Sanskar Chaudhary**
-   **Reika Yonetsu**

* * * * *

**Testing**
-----------

### **Unit and Integration Testing**

The application includes unit tests to ensure its stability and functionality:

-   Unit tests for individual components.

To run tests:

`go test ./...`

* * * * *

**Contact and Support**
-----------------------

For any inquiries or support, please email us at: support@invoicerator.com.

* * * * *
