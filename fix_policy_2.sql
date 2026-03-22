-- Completely drop the problematic policies
DROP POLICY IF EXISTS "Users can view members of their teams" ON public.team_members;
DROP POLICY IF EXISTS "Users can view teams they belong to" ON public.teams;
DROP POLICY IF EXISTS "Users can update teams they belong to" ON public.teams;

-- Recreate them safely using a simpler approach
CREATE POLICY "Users can view their own team memberships" 
  ON public.team_members FOR SELECT 
  USING (user_id = auth.uid());

CREATE POLICY "Users can view teams they belong to" 
  ON public.teams FOR SELECT 
  USING (id IN (SELECT team_id FROM public.team_members WHERE user_id = auth.uid()));

CREATE POLICY "Users can update teams they belong to" 
  ON public.teams FOR UPDATE 
  USING (id IN (SELECT team_id FROM public.team_members WHERE user_id = auth.uid()));
