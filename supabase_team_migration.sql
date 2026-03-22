-- 1. Create teams table
CREATE TABLE public.teams (
  id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  company_name TEXT,
  company_email TEXT,
  company_address TEXT,
  company_phone TEXT,
  letterhead_url TEXT,
  selected_fields JSONB DEFAULT '[]'::jsonb,
  terms_conditions TEXT,
  invite_code UUID DEFAULT gen_random_uuid() UNIQUE,
  created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 2. Create team_members table
CREATE TABLE public.team_members (
  id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  team_id UUID REFERENCES public.teams(id) ON DELETE CASCADE NOT NULL,
  user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(team_id, user_id)
);

-- 3. Migrate existing data from profiles to teams
DO $$
DECLARE
  prof RECORD;
  new_team_id UUID;
BEGIN
  FOR prof IN SELECT * FROM public.profiles LOOP
    INSERT INTO public.teams (company_name, company_email, company_address, company_phone, letterhead_url, selected_fields, terms_conditions, created_by)
    VALUES (prof.company_name, prof.company_email, prof.company_address, prof.company_phone, prof.letterhead_url, prof.selected_fields, prof.terms_conditions, prof.id)
    RETURNING id INTO new_team_id;

    INSERT INTO public.team_members (team_id, user_id)
    VALUES (new_team_id, prof.id);
  END LOOP;
END $$;

-- 4. Enable RLS
ALTER TABLE public.teams ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.team_members ENABLE ROW LEVEL SECURITY;

-- 5. RLS Policies for teams
CREATE POLICY "Users can view teams they belong to" 
  ON public.teams FOR SELECT 
  USING (EXISTS (SELECT 1 FROM public.team_members WHERE team_id = id AND user_id = auth.uid()));

CREATE POLICY "Users can update teams they belong to" 
  ON public.teams FOR UPDATE 
  USING (EXISTS (SELECT 1 FROM public.team_members WHERE team_id = id AND user_id = auth.uid()));

-- 6. RLS Policies for team_members
CREATE POLICY "Users can view members of their teams" 
  ON public.team_members FOR SELECT 
  USING (
    user_id = auth.uid() OR 
    EXISTS (SELECT 1 FROM public.team_members tm WHERE tm.team_id = team_id AND tm.user_id = auth.uid())
  );

-- 7. Secure function to join a team using an invite code
CREATE OR REPLACE FUNCTION public.join_team_by_code(invite_code_param UUID)
RETURNS BOOLEAN AS $$
DECLARE
  target_team_id UUID;
BEGIN
  -- Find the team with the matching invite code
  SELECT id INTO target_team_id FROM public.teams WHERE invite_code = invite_code_param;
  
  -- If no team found, return false
  IF target_team_id IS NULL THEN
    RETURN FALSE;
  END IF;

  -- Add the user to the team (ignore if they are already in it)
  INSERT INTO public.team_members (team_id, user_id)
  VALUES (target_team_id, auth.uid())
  ON CONFLICT DO NOTHING;
  
  RETURN TRUE;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- 8. Update the handle_new_user trigger to create a team instead of a profile
CREATE OR REPLACE FUNCTION public.handle_new_user() 
RETURNS TRIGGER AS $$
DECLARE
  new_team_id UUID;
BEGIN
  -- Create a default team for the new user (leave company name blank)
  INSERT INTO public.teams (company_name, created_by)
  VALUES (NULL, new.id)
  RETURNING id INTO new_team_id;

  -- Add the user to their new team
  INSERT INTO public.team_members (team_id, user_id)
  VALUES (new_team_id, new.id);

  RETURN new;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Note: We are leaving the old profiles table intact for now as a backup, 
-- but the app will no longer use it.
