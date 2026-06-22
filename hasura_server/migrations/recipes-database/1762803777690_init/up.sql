SET transaction_timeout = 0;
SET check_function_bodies = false;
CREATE FUNCTION public.set_current_timestamp_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
  _new record;
BEGIN
  _new := NEW;
  _new."updated_at" = NOW();
  RETURN _new;
END;
$$;
CREATE FUNCTION public.set_expire_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW."expireAt" := NOW() + INTERVAL '15 minutes';
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_recipe_avg_rating() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    UPDATE "Recipes"
    SET avg_rating = (
        SELECT ROUND(AVG(rating), 2)
        FROM "Ratings"
        WHERE recipe_id = NEW.recipe_id
    )
    WHERE id = NEW.recipe_id;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_recipe_like_count() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    UPDATE "Recipes"
    SET like_count = (
        SELECT COUNT(*)
        FROM "Likes"
        WHERE recipe_id = COALESCE(NEW.recipe_id, OLD.recipe_id)
    )
    WHERE id = COALESCE(NEW.recipe_id, OLD.recipe_id);
    RETURN NEW;
END;
$$;
CREATE TABLE public."Bookmarks" (
    user_id uuid NOT NULL,
    recipe_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);
COMMENT ON TABLE public."Bookmarks" IS 'users favorite recipes are store here.  ';
CREATE TABLE public."Categories" (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name text NOT NULL,
    image_url text NOT NULL,
    description text,
    created_at timestamp with time zone DEFAULT now()
);
COMMENT ON TABLE public."Categories" IS 'different categories for those recipe are list here in this relation.';
CREATE TABLE public."Comments" (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    recipe_id uuid NOT NULL,
    user_id uuid NOT NULL,
    text text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);
COMMENT ON TABLE public."Comments" IS 'All users comment on a specific recipe store here.';
CREATE TABLE public."Ingredient" (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    recipe_id uuid NOT NULL,
    name text NOT NULL,
    quantity text NOT NULL
);
COMMENT ON TABLE public."Ingredient" IS 'All ingredient for each recipe are listed here.';
CREATE TABLE public."Likes" (
    user_id uuid NOT NULL,
    recipe_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);
COMMENT ON TABLE public."Likes" IS 'All likes for each recipe by logged in user are store here.';
CREATE TABLE public."Ratings" (
    user_id uuid NOT NULL,
    recipe_id uuid NOT NULL,
    rating double precision NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);
COMMENT ON TABLE public."Ratings" IS 'All user rating for each recipe are store here.';
CREATE TABLE public."Recipes" (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    title text NOT NULL,
    description text NOT NULL,
    created_id uuid NOT NULL,
    category_id uuid NOT NULL,
    prep_time_minutes integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    avg_rating double precision DEFAULT '0'::double precision,
    like_count integer DEFAULT 0
);
COMMENT ON TABLE public."Recipes" IS 'All Recipes information of the site created by user store here.';
CREATE TABLE public."Recipes_image" (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    recipe_id uuid NOT NULL,
    image_url text NOT NULL,
    is_featured boolean DEFAULT false NOT NULL
);
COMMENT ON TABLE public."Recipes_image" IS 'All the recipes thumbnail image url store here including the featured image.';
CREATE TABLE public."Users" (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name text NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    avater_url text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    "isVerified" boolean DEFAULT false NOT NULL
);
COMMENT ON TABLE public."Users" IS 'All website users basic information store here';
CREATE TABLE public."VerificationData" (
    email text NOT NULL,
    code text NOT NULL,
    type text NOT NULL,
    "expireAt" timestamp with time zone
);
COMMENT ON TABLE public."VerificationData" IS 'Used to store user email and related thing till they verify it.';
ALTER TABLE ONLY public."Bookmarks"
    ADD CONSTRAINT "Bookmarks_pkey" PRIMARY KEY (user_id, recipe_id);
ALTER TABLE ONLY public."Categories"
    ADD CONSTRAINT "Categories_pkey" PRIMARY KEY (id);
ALTER TABLE ONLY public."Comments"
    ADD CONSTRAINT "Comments_pkey" PRIMARY KEY (id);
ALTER TABLE ONLY public."Ingredient"
    ADD CONSTRAINT "Ingredient_pkey" PRIMARY KEY (id);
ALTER TABLE ONLY public."Likes"
    ADD CONSTRAINT "Likes_pkey" PRIMARY KEY (user_id, recipe_id);
ALTER TABLE ONLY public."Ratings"
    ADD CONSTRAINT "Ratings_pkey" PRIMARY KEY (recipe_id, user_id);
ALTER TABLE ONLY public."Recipes_image"
    ADD CONSTRAINT "Recipes_image_pkey" PRIMARY KEY (id);
ALTER TABLE ONLY public."Recipes"
    ADD CONSTRAINT "Recipes_pkey" PRIMARY KEY (id);
ALTER TABLE ONLY public."Users"
    ADD CONSTRAINT "Users_email_key" UNIQUE (email);
ALTER TABLE ONLY public."Users"
    ADD CONSTRAINT "Users_pkey" PRIMARY KEY (id);
ALTER TABLE ONLY public."VerificationData"
    ADD CONSTRAINT "VerificationData_pkey" PRIMARY KEY (email);
CREATE TRIGGER "set_public_Recipes_updated_at" BEFORE UPDATE ON public."Recipes" FOR EACH ROW EXECUTE FUNCTION public.set_current_timestamp_updated_at();
COMMENT ON TRIGGER "set_public_Recipes_updated_at" ON public."Recipes" IS 'trigger to set value of column "updated_at" to current timestamp on row update';
CREATE TRIGGER "set_public_Users_updated_at" BEFORE UPDATE ON public."Users" FOR EACH ROW EXECUTE FUNCTION public.set_current_timestamp_updated_at();
COMMENT ON TRIGGER "set_public_Users_updated_at" ON public."Users" IS 'trigger to set value of column "updated_at" to current timestamp on row update';
CREATE TRIGGER trg_update_avg_rating AFTER INSERT OR DELETE OR UPDATE ON public."Ratings" FOR EACH ROW EXECUTE FUNCTION public.update_recipe_avg_rating();
CREATE TRIGGER trg_update_like_count AFTER INSERT OR DELETE ON public."Likes" FOR EACH ROW EXECUTE FUNCTION public.update_recipe_like_count();
CREATE TRIGGER trigger_set_expire_at BEFORE INSERT ON public."VerificationData" FOR EACH ROW EXECUTE FUNCTION public.set_expire_at();
ALTER TABLE ONLY public."Bookmarks"
    ADD CONSTRAINT "Bookmarks_recipe_id_fkey" FOREIGN KEY (recipe_id) REFERENCES public."Recipes"(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public."Bookmarks"
    ADD CONSTRAINT "Bookmarks_user_id_fkey" FOREIGN KEY (user_id) REFERENCES public."Users"(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public."Comments"
    ADD CONSTRAINT "Comments_recipe_id_fkey" FOREIGN KEY (recipe_id) REFERENCES public."Recipes"(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public."Comments"
    ADD CONSTRAINT "Comments_user_id_fkey" FOREIGN KEY (user_id) REFERENCES public."Users"(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public."Ingredient"
    ADD CONSTRAINT "Ingredient_recipe_id_fkey" FOREIGN KEY (recipe_id) REFERENCES public."Recipes"(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public."Likes"
    ADD CONSTRAINT "Likes_recipe_id_fkey" FOREIGN KEY (recipe_id) REFERENCES public."Recipes"(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public."Likes"
    ADD CONSTRAINT "Likes_user_id_fkey" FOREIGN KEY (user_id) REFERENCES public."Users"(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public."Ratings"
    ADD CONSTRAINT "Ratings_recipe_id_fkey" FOREIGN KEY (recipe_id) REFERENCES public."Recipes"(id) ON UPDATE RESTRICT ON DELETE CASCADE;
ALTER TABLE ONLY public."Ratings"
    ADD CONSTRAINT "Ratings_user_id_fkey" FOREIGN KEY (user_id) REFERENCES public."Users"(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public."Recipes"
    ADD CONSTRAINT "Recipes_category_id_fkey" FOREIGN KEY (category_id) REFERENCES public."Categories"(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public."Recipes"
    ADD CONSTRAINT "Recipes_created_id_fkey" FOREIGN KEY (created_id) REFERENCES public."Users"(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public."Recipes_image"
    ADD CONSTRAINT "Recipes_image_recipe_id_fkey" FOREIGN KEY (recipe_id) REFERENCES public."Recipes"(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public."VerificationData"
    ADD CONSTRAINT "VerificationData_email_fkey" FOREIGN KEY (email) REFERENCES public."Users"(email) ON UPDATE CASCADE ON DELETE CASCADE;

