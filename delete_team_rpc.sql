-- Team lead only. Deletes the team row; ON DELETE CASCADE removes team_members and team_invites.
CREATE OR REPLACE FUNCTION public.delete_team_for_lead(p_team_id UUID)
RETURNS void
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = public
AS $$
DECLARE
  lead_id UUID;
  caller_id UUID;
  mcount int;
BEGIN
  caller_id := coalesce(
    auth.uid(),
    (SELECT NULLIF(trim(auth.jwt() ->> 'sub'), '')::uuid)
  );

  IF caller_id IS NULL THEN
    RAISE EXCEPTION 'Not authenticated';
  END IF;

  IF NOT EXISTS (
    SELECT 1 FROM public.team_members WHERE team_id = p_team_id AND user_id = caller_id
  ) THEN
    RAISE EXCEPTION 'Not a member of this team';
  END IF;

  SELECT created_by INTO lead_id FROM public.teams WHERE id = p_team_id;

  IF lead_id IS NOT NULL AND lead_id <> caller_id THEN
    RAISE EXCEPTION 'Only the team lead can delete this team';
  END IF;

  -- Legacy rows with no created_by: only the sole member may delete
  IF lead_id IS NULL THEN
    SELECT COUNT(*)::int INTO mcount FROM public.team_members WHERE team_id = p_team_id;
    IF mcount <> 1 OR NOT EXISTS (
      SELECT 1 FROM public.team_members WHERE team_id = p_team_id AND user_id = caller_id
    ) THEN
      RAISE EXCEPTION 'Only the team lead can delete this team';
    END IF;
  END IF;

  DELETE FROM public.teams WHERE id = p_team_id;
END;
$$;

GRANT EXECUTE ON FUNCTION public.delete_team_for_lead(UUID) TO authenticated;
