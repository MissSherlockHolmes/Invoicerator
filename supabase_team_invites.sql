-- 1. Create the team_invites table
CREATE TABLE public.team_invites (
  id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  team_id UUID REFERENCES public.teams(id) ON DELETE CASCADE NOT NULL,
  email TEXT NOT NULL,
  token UUID DEFAULT gen_random_uuid() UNIQUE NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(team_id, email)
);

-- 2. Enable RLS
ALTER TABLE public.team_invites ENABLE ROW LEVEL SECURITY;

-- 3. Policy: Team members can manage invites for their team
CREATE POLICY "Team members can manage invites" 
  ON public.team_invites FOR ALL 
  USING (team_id IN (SELECT team_id FROM public.team_members WHERE user_id = auth.uid()));

-- 4. Secure function to accept an invite
CREATE OR REPLACE FUNCTION public.accept_team_invite(invite_token UUID)
RETURNS BOOLEAN AS $$
DECLARE
  target_invite RECORD;
  user_email TEXT;
BEGIN
  -- Get the user's email from auth.users
  SELECT email INTO user_email FROM auth.users WHERE id = auth.uid();

  -- Find the invite
  SELECT * INTO target_invite FROM public.team_invites WHERE token = invite_token;
  
  IF target_invite IS NULL THEN
    RAISE EXCEPTION 'Invalid invite token.';
  END IF;

  -- Check if the logged-in user's email matches the invited email
  IF lower(target_invite.email) != lower(user_email) THEN
    RAISE EXCEPTION 'This invite was sent to a different email address. Please log in with %', target_invite.email;
  END IF;

  -- Add user to team
  INSERT INTO public.team_members (team_id, user_id)
  VALUES (target_invite.team_id, auth.uid())
  ON CONFLICT DO NOTHING;

  -- Delete the invite so it can't be used again
  DELETE FROM public.team_invites WHERE id = target_invite.id;

  RETURN TRUE;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- 5. (Optional) Remove the old generic invite code from teams
-- ALTER TABLE public.teams DROP COLUMN IF EXISTS invite_code;
-- DROP FUNCTION IF EXISTS public.join_team_by_code(UUID);
