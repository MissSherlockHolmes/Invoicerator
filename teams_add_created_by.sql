-- Run once on Supabase: track who created each team (team lead).
-- Only the team lead may remove members (enforced in remove_team_member).

ALTER TABLE public.teams
  ADD COLUMN IF NOT EXISTS created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL;

-- Backfill legacy teams: treat the earliest joined member as the creator.
UPDATE public.teams t
SET created_by = s.user_id
FROM (
  SELECT DISTINCT ON (team_id) team_id, user_id
  FROM public.team_members
  ORDER BY team_id, created_at ASC NULLS LAST, user_id ASC
) s
WHERE t.id = s.team_id
  AND (t.created_by IS NULL);

COMMENT ON COLUMN public.teams.created_by IS 'User who created the team; only this user may remove members or manage invites (app-enforced).';
