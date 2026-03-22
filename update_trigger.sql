CREATE OR REPLACE FUNCTION public.handle_new_user() 
RETURNS TRIGGER AS $$
DECLARE
  new_team_id UUID;
BEGIN
  -- Create a default team for the new user (leave company name blank)
  INSERT INTO public.teams (company_name)
  VALUES (NULL)
  RETURNING id INTO new_team_id;

  -- Add the user to their new team
  INSERT INTO public.team_members (team_id, user_id)
  VALUES (new_team_id, new.id);

  RETURN new;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;
