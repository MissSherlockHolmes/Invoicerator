-- Supabase Schema for Invoicerator

-- 1. Profiles Table (extends Supabase auth.users)
CREATE TABLE public.profiles (
  id UUID REFERENCES auth.users(id) ON DELETE CASCADE PRIMARY KEY,
  username TEXT UNIQUE,
  company_name TEXT,
  company_email TEXT,
  company_address TEXT,
  company_phone TEXT,
  letterhead_url TEXT,
  selected_fields JSONB DEFAULT '[]'::jsonb,
  terms_conditions TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Enable Row Level Security (RLS)
ALTER TABLE public.profiles ENABLE ROW LEVEL SECURITY;

-- Create policies for profiles
CREATE POLICY "Users can view their own profile" 
  ON public.profiles FOR SELECT 
  USING (auth.uid() = id);

CREATE POLICY "Users can update their own profile" 
  ON public.profiles FOR UPDATE 
  USING (auth.uid() = id);

CREATE POLICY "Users can insert their own profile" 
  ON public.profiles FOR INSERT 
  WITH CHECK (auth.uid() = id);

-- 2. Financial Institutions Table
CREATE TABLE public.financial_institutions (
  id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  user_id UUID REFERENCES public.profiles(id) ON DELETE CASCADE NOT NULL,
  name TEXT NOT NULL,
  bank_number TEXT,
  link TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Enable Row Level Security (RLS)
ALTER TABLE public.financial_institutions ENABLE ROW LEVEL SECURITY;

-- Create policies for financial institutions
CREATE POLICY "Users can view their own financial institutions" 
  ON public.financial_institutions FOR SELECT 
  USING (auth.uid() = user_id);

CREATE POLICY "Users can insert their own financial institutions" 
  ON public.financial_institutions FOR INSERT 
  WITH CHECK (auth.uid() = user_id);

CREATE POLICY "Users can update their own financial institutions" 
  ON public.financial_institutions FOR UPDATE 
  USING (auth.uid() = user_id);

CREATE POLICY "Users can delete their own financial institutions" 
  ON public.financial_institutions FOR DELETE 
  USING (auth.uid() = user_id);

-- 3. Storage Bucket for Letterheads
-- Note: Create a bucket named 'letterheads' in the Supabase Storage dashboard.
-- Set it to 'public' if you want the images to be publicly accessible, 
-- or 'private' and use signed URLs.

-- Example Storage Policies (if creating via SQL):
-- INSERT INTO storage.buckets (id, name, public) VALUES ('letterheads', 'letterheads', true);
-- CREATE POLICY "Users can upload their own letterhead"
--   ON storage.objects FOR INSERT
--   WITH CHECK (bucket_id = 'letterheads' AND auth.uid()::text = (storage.foldername(name))[1]);
-- CREATE POLICY "Anyone can view letterheads"
--   ON storage.objects FOR SELECT
--   USING (bucket_id = 'letterheads');

-- 4. Trigger to automatically create a profile when a new user signs up
CREATE OR REPLACE FUNCTION public.handle_new_user() 
RETURNS TRIGGER AS $$
BEGIN
  INSERT INTO public.profiles (id, username)
  VALUES (new.id, new.raw_user_meta_data->>'username');
  RETURN new;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

CREATE TRIGGER on_auth_user_created
  AFTER INSERT ON auth.users
  FOR EACH ROW EXECUTE FUNCTION public.handle_new_user();
