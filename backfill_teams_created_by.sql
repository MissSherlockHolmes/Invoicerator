-- Run in Supabase SQL editor if teams.created_by is NULL for some rows.
-- Assigns creator to the earliest team_members row per team (matches teams_add_created_by.sql).

UPDATE public.teams t
SET created_by = s.user_id
FROM (
  SELECT DISTINCT ON (team_id) team_id, user_id
  FROM public.team_members
  ORDER BY team_id, created_at ASC NULLS LAST, user_id ASC
) s
WHERE t.id = s.team_id
  AND t.created_by IS NULL;
