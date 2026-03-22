-- Drop the existing policy that's too restrictive
DROP POLICY IF EXISTS "Users can insert teams" ON public.teams;

-- Create a new policy that allows any authenticated user to create a team
CREATE POLICY "Users can insert teams" 
ON public.teams FOR INSERT 
TO authenticated 
WITH CHECK (true);

-- Ensure team_members also has an insert policy so they can add themselves to the team they just created
DROP POLICY IF EXISTS "Users can insert team_members" ON public.team_members;

CREATE POLICY "Users can insert team_members" 
ON public.team_members FOR INSERT 
TO authenticated 
WITH CHECK (user_id = auth.uid());
