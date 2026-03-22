CREATE OR REPLACE FUNCTION public.create_new_team(new_company_name TEXT)
RETURNS UUID
LANGUAGE plpgsql
SECURITY DEFINER
AS $$
DECLARE
  new_team_id UUID;
BEGIN
  -- Insert the new team
  INSERT INTO public.teams (company_name)
  VALUES (new_company_name)
  RETURNING id INTO new_team_id;

  -- Add the calling user to the new team
  INSERT INTO public.team_members (team_id, user_id)
  VALUES (new_team_id, auth.uid());

  RETURN new_team_id;
END;
$$;
