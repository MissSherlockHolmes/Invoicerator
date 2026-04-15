// public/js/supabase.js
import { createClient } from 'https://cdn.jsdelivr.net/npm/@supabase/supabase-js@2/+esm'

// Replace these with your actual Supabase project URL and anon key
const supabaseUrl = 'https://blqzmaxcbsegeskvedwv.supabase.co'
const supabaseAnonKey = 'sb_publishable_2_NDf00PPwBgzM8uH63I1g_byoOVg31'

let supabase = null
try {
    supabase = createClient(supabaseUrl, supabaseAnonKey)
} catch (err) {
    console.error('Invoicerator: failed to create Supabase client', err)
}

export { supabase }

export function isSupabaseReady() {
    return supabase != null
}
