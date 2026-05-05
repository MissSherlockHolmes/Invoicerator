-- Keep-alive heartbeat primitives for reliable Supabase activity verification.
-- Run this in Supabase SQL Editor once.

CREATE TABLE IF NOT EXISTS public.keepalive_heartbeat (
  id INTEGER PRIMARY KEY DEFAULT 1 CHECK (id = 1),
  source TEXT NOT NULL DEFAULT 'unknown',
  last_seen_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION public.keepalive_heartbeat_touch(p_source TEXT DEFAULT 'github-actions')
RETURNS TABLE(last_seen_at TIMESTAMPTZ, source TEXT)
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = public
AS $$
BEGIN
  INSERT INTO public.keepalive_heartbeat (id, source, last_seen_at, updated_at)
  VALUES (1, COALESCE(NULLIF(TRIM(p_source), ''), 'unknown'), NOW(), NOW())
  ON CONFLICT (id) DO UPDATE
  SET source = EXCLUDED.source,
      last_seen_at = NOW(),
      updated_at = NOW();

  RETURN QUERY
  SELECT k.last_seen_at, k.source
  FROM public.keepalive_heartbeat k
  WHERE k.id = 1;
END;
$$;

REVOKE ALL ON TABLE public.keepalive_heartbeat FROM PUBLIC;
GRANT SELECT ON TABLE public.keepalive_heartbeat TO service_role;
GRANT EXECUTE ON FUNCTION public.keepalive_heartbeat_touch(TEXT) TO service_role;
