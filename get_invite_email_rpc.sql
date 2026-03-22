CREATE OR REPLACE FUNCTION public.get_invite_email_by_token(p_token UUID)
RETURNS TEXT
LANGUAGE plpgsql
SECURITY DEFINER
AS $$
DECLARE
  v_email TEXT;
BEGIN
  SELECT email INTO v_email
  FROM public.team_invites
  WHERE token = p_token;
  
  RETURN v_email;
END;
$$;
