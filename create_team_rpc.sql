-- Caller is always stored as teams.created_by (team lead).
-- Uses JWT sub as fallback if auth.uid() is ever null in SECURITY DEFINER context.
CREATE OR REPLACE FUNCTION public.create_new_team(new_company_name TEXT)
RETURNS UUID
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = public
AS $$
DECLARE
  new_team_id UUID;
  caller_id UUID;
BEGIN
  caller_id := coalesce(
    auth.uid(),
    (SELECT NULLIF(trim(auth.jwt() ->> 'sub'), '')::uuid)
  );

  IF caller_id IS NULL THEN
    RAISE EXCEPTION 'Not authenticated';
  END IF;

  INSERT INTO public.teams (company_name, created_by)
  VALUES (new_company_name, caller_id)
  RETURNING id INTO new_team_id;

  INSERT INTO public.team_members (team_id, user_id)
  VALUES (new_team_id, caller_id);

  RETURN new_team_id;
END;
$$;

GRANT EXECUTE ON FUNCTION public.create_new_team(TEXT) TO authenticated;
