// public/js/supabase.js
import { createClient } from 'https://cdn.jsdelivr.net/npm/@supabase/supabase-js@2/+esm'

// Replace these with your actual Supabase project URL and anon key
const supabaseUrl = 'https://blqzmaxcbsegeskvedwv.supabase.co'
const supabaseAnonKey = 'sb_publishable_2_NDf00PPwBgzM8uH63I1g_byoOVg31'

export const supabase = createClient(supabaseUrl, supabaseAnonKey)
