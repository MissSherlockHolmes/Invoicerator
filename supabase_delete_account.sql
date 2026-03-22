-- Function to allow a user to delete their own account
-- This is necessary because the client-side Supabase JS library cannot delete auth.users directly for security reasons.
CREATE OR REPLACE FUNCTION public.delete_user_account()
RETURNS void AS $$
BEGIN
  -- Delete the user from the auth.users table.
  -- Because we set up ON DELETE CASCADE on the team_members table, 
  -- this will automatically remove them from their teams as well.
  DELETE FROM auth.users WHERE id = auth.uid();
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;
