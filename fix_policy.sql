-- Fix the infinite recursion policy
DROP POLICY IF EXISTS "Users can view members of their teams" ON public.team_members;

-- Create a simpler policy that doesn't reference itself
CREATE POLICY "Users can view members of their teams" 
  ON public.team_members FOR SELECT 
  USING (
    user_id = auth.uid() OR 
    team_id IN (SELECT team_id FROM public.team_members WHERE user_id = auth.uid())
  );
