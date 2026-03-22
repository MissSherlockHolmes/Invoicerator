CREATE OR REPLACE FUNCTION public.handle_new_user() 
RETURNS TRIGGER AS $$
DECLARE
  new_team_id UUID;
  invite_record RECORD;
  has_invites BOOLEAN := FALSE;
BEGIN
  -- 1. Find all pending invites for this user's email
  FOR invite_record IN SELECT * FROM public.team_invites WHERE email = new.email LOOP
    has_invites := TRUE;
    
    -- Add the user to the team they were invited to
    INSERT INTO public.team_members (team_id, user_id)
    VALUES (invite_record.team_id, new.id)
    ON CONFLICT DO NOTHING;
    
    -- Delete the used invite
    DELETE FROM public.team_invites WHERE token = invite_record.token;
  END LOOP;

  -- 2. If they had no invites, create a default blank team for them
  IF NOT has_invites THEN
    INSERT INTO public.teams (company_name)
    VALUES (NULL)
    RETURNING id INTO new_team_id;

    INSERT INTO public.team_members (team_id, user_id)
    VALUES (new_team_id, new.id);
  END IF;

  RETURN new;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;
