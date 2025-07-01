CREATE SCHEMA IF NOT EXISTS auth;

ALTER TABLE public.users SET SCHEMA auth;
ALTER TABLE public.refresh_tokens SET SCHEMA auth;