import { supabase } from './supabase.js';

export async function renderNavbar() {
    const container = document.getElementById('navbar-container');
    if (!container) return;

    const { data: { session } } = await supabase.auth.getSession();

    let rightLinks = '';
    if (session) {
        rightLinks = `
            <li class="nav-item">
                <a class="nav-link" href="/profile.html">Profile</a>
            </li>
            <li class="nav-item">
                <a class="nav-link" href="#" id="logout-btn">Logout</a>
            </li>
        `;
    } else {
        rightLinks = `
            <li class="nav-item">
                <a class="nav-link" href="/login.html">Login</a>
            </li>
            <li class="nav-item">
                <a class="nav-link" href="/signup.html">Sign Up</a>
            </li>
        `;
    }

    container.innerHTML = `
        <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
            <div class="container-fluid">
                <a class="navbar-brand" href="/">Invoicerator</a>
                <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav" aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
                    <span class="navbar-toggler-icon"></span>
                </button>
                <div class="collapse navbar-collapse" id="navbarNav">
                    <ul class="navbar-nav me-auto">
                        <li class="nav-item">
                            <a class="nav-link" href="/options.html">Options</a>
                        </li>
                    </ul>
                    <ul class="navbar-nav">
                        ${rightLinks}
                    </ul>
                </div>
            </div>
        </nav>
    `;

    if (session) {
        document.getElementById('logout-btn').addEventListener('click', async (e) => {
            e.preventDefault();
            await supabase.auth.signOut();
            window.location.href = '/';
        });
    }
}
