-- 1. Function to get team members with their email addresses securely
CREATE OR REPLACE FUNCTION get_team_members_with_emails(p_team_id UUID)
RETURNS TABLE (user_id UUID, email TEXT)
LANGUAGE plpgsql
SECURITY DEFINER
AS $$
BEGIN
  -- Check if the calling user is a member of this team
  IF NOT EXISTS (SELECT 1 FROM team_members WHERE team_id = p_team_id AND team_members.user_id = auth.uid()) THEN
    RAISE EXCEPTION 'Not authorized';
  END IF;

  RETURN QUERY
  SELECT tm.user_id, au.email::TEXT
  FROM team_members tm
  JOIN auth.users au ON tm.user_id = au.id
  WHERE tm.team_id = p_team_id;
END;
$$;

-- 2. Function to remove a team member securely (team lead only)
CREATE OR REPLACE FUNCTION remove_team_member(p_team_id UUID, p_user_id UUID)
RETURNS BOOLEAN
LANGUAGE plpgsql
SECURITY DEFINER
AS $$
DECLARE
  lead_id UUID;
BEGIN
  -- Caller must be a member
  IF NOT EXISTS (SELECT 1 FROM team_members WHERE team_id = p_team_id AND team_members.user_id = auth.uid()) THEN
    RAISE EXCEPTION 'Not authorized';
  END IF;

  SELECT created_by INTO lead_id FROM public.teams WHERE id = p_team_id;

  -- Only the user who created the team may remove others
  IF lead_id IS NULL OR auth.uid() <> lead_id THEN
    RAISE EXCEPTION 'Only the team lead can remove members';
  END IF;

  -- Cannot remove the team lead
  IF p_user_id = lead_id THEN
    RAISE EXCEPTION 'Cannot remove the team lead';
  END IF;

  -- Prevent removing the last member of the team (optional, but good practice)
  IF (SELECT COUNT(*) FROM team_members WHERE team_id = p_team_id) <= 1 THEN
    RAISE EXCEPTION 'Cannot remove the only member of the team. Delete the team instead.';
  END IF;

  DELETE FROM team_members WHERE team_id = p_team_id AND user_id = p_user_id;
  RETURN TRUE;
END;
$$;
