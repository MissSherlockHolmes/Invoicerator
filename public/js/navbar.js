import { supabase } from './supabase.js';
import { getSessionSafe } from './supabase-connection.js';

export async function renderNavbar() {
    const container = document.getElementById('navbar-container');
    if (!container) return { ok: true, session: null };

    const sessionResult = await getSessionSafe(supabase);
    if (!sessionResult.ok) {
        container.innerHTML = `
        <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
            <div class="container-fluid">
                <span class="navbar-brand">Invoicerator</span>
                <span class="navbar-text text-warning small">Internal connection error</span>
            </div>
        </nav>`;
        return { ok: false, session: null };
    }

    const session = sessionResult.data?.session ?? null;
    const path = window.location.pathname;

    const homeLink = session ? '/options.html' : '/';

    const escapeHtml = (s) => String(s)
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;');

    let rightLinks = '';

    if (session) {
        const email = session.user?.email || '';
        rightLinks += `<li class="nav-item d-flex align-items-center me-lg-2 mb-2 mb-lg-0">
            <span class="navbar-text text-secondary small text-truncate" style="max-width: 14rem;" title="${escapeHtml(email)}">Signed in as <span class="text-light">${escapeHtml(email)}</span></span>
        </li>`;
        // If we are NOT on the welcome/options page, show links to the other pages
        if (path !== '/options.html') {
            if (path !== '/create_invoice.html') {
                rightLinks += `<li class="nav-item"><a class="nav-link" href="/create_invoice.html">Create Invoice</a></li>`;
            }
            if (path !== '/edit_invoice.html') {
                rightLinks += `<li class="nav-item"><a class="nav-link" href="/edit_invoice.html">Edit Invoice</a></li>`;
            }
            if (path !== '/profile.html') {
                rightLinks += `<li class="nav-item"><a class="nav-link" href="/profile.html">Profile</a></li>`;
            }
        }
        // Always show Logout
        rightLinks += `<li class="nav-item"><a class="nav-link text-danger" href="#" id="logout-btn">Logout</a></li>`;
    } else {
        // Unauthenticated links
        if (path !== '/login.html') {
            rightLinks += `<li class="nav-item"><a class="nav-link" href="/login.html">Login</a></li>`;
        }
        if (path !== '/signup.html') {
            rightLinks += `<li class="nav-item"><a class="nav-link" href="/signup.html">Sign Up</a></li>`;
        }
    }

    container.innerHTML = `
        <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
            <div class="container-fluid">
                <a class="navbar-brand d-flex align-items-center" href="${homeLink}" title="Home">
                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="currentColor" class="bi bi-house-door-fill me-2" viewBox="0 0 16 16">
                      <path d="M6.5 14.5v-3.505c0-.245.25-.495.5-.495h2c.25 0 .5.25.5.5v3.5a.5.5 0 0 0 .5.5h4a.5.5 0 0 0 .5-.5v-7a.5.5 0 0 0-.146-.354L13 5.793V2.5a.5.5 0 0 0-.5-.5h-1a.5.5 0 0 0-.5.5v1.293L8.354 1.146a.5.5 0 0 0-.708 0l-6 6A.5.5 0 0 0 1.5 7.5v7a.5.5 0 0 0 .5.5h4a.5.5 0 0 0 .5-.5z"/>
                    </svg>
                    Invoicerator
                </a>
                <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav" aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
                    <span class="navbar-toggler-icon"></span>
                </button>
                <div class="collapse navbar-collapse" id="navbarNav">
                    <ul class="navbar-nav ms-auto">
                        ${rightLinks}
                    </ul>
                </div>
            </div>
        </nav>
    `;

    if (session) {
        const logoutBtn = document.getElementById('logout-btn');
        if (logoutBtn) {
            logoutBtn.addEventListener('click', async (e) => {
                e.preventDefault();
                await supabase.auth.signOut();
                window.location.href = '/';
            });
        }
    }

    return { ok: true, session };
}
