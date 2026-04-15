/**
 * Safe session read: catches DNS/network failures when the Supabase project
 * is paused or unreachable (fetch throws before a structured Auth error).
 */
export async function getSessionSafe(supabaseClient) {
    if (!supabaseClient) {
        return { ok: false, error: new Error('Supabase client is not initialized') };
    }
    try {
        const { data, error } = await supabaseClient.auth.getSession();
        return { ok: true, data, error };
    } catch (error) {
        return { ok: false, error };
    }
}

/** Friendly banner when the app cannot reach the server (no infra jargon for end users). */
export function serviceUnavailableMarkupDetailed() {
    return `
        <div class="alert alert-warning text-start mx-auto" style="max-width: 28rem;" role="alert">
            <p class="mb-3">Internal connection error, please contact the site administrator or check back in an hour.</p>
            <p class="mb-0 small"><a href="/" class="alert-link">Back to home</a></p>
        </div>
    `;
}
