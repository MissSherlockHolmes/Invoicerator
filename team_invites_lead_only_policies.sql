-- Restrict team_invites mutations to the team creator (teams.created_by).
-- All members can still SELECT invites for teams they belong to (see pending list).

DROP POLICY IF EXISTS "Team members can manage invites" ON public.team_invites;

-- Anyone on the team can see pending invites (read-only for non-leads in the app)
CREATE POLICY "Team members can view invites for their teams"
  ON public.team_invites FOR SELECT
  USING (
    team_id IN (SELECT team_id FROM public.team_members WHERE user_id = auth.uid())
  );

CREATE POLICY "Team leads can insert invites"
  ON public.team_invites FOR INSERT
  WITH CHECK (
    EXISTS (
      SELECT 1 FROM public.teams t
      WHERE t.id = team_id AND t.created_by IS NOT NULL AND t.created_by = auth.uid()
    )
  );

CREATE POLICY "Team leads can update invites"
  ON public.team_invites FOR UPDATE
  USING (
    EXISTS (
      SELECT 1 FROM public.teams t
      WHERE t.id = team_id AND t.created_by IS NOT NULL AND t.created_by = auth.uid()
    )
  );

CREATE POLICY "Team leads can delete invites"
  ON public.team_invites FOR DELETE
  USING (
    EXISTS (
      SELECT 1 FROM public.teams t
      WHERE t.id = team_id AND t.created_by IS NOT NULL AND t.created_by = auth.uid()
    )
  );
