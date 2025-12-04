-- Migration: Add Full-Text Search support to products table
-- This migration adds a GIN index on the name column for full-text search

-- 1. Add a tsvector column to store the search vector (optional but recommended for performance)
ALTER TABLE products ADD COLUMN IF NOT EXISTS name_search tsvector
    GENERATED ALWAYS AS (to_tsvector('english', COALESCE(name, ''))) STORED;

-- 2. Create a GIN index on the tsvector column for fast full-text search
CREATE INDEX IF NOT EXISTS products_name_search_idx ON products USING GIN (name_search);

-- Alternative approach (if you don't want to add a new column):
-- You can create a GIN index directly on the to_tsvector expression
-- Uncomment the line below if you prefer this approach (and comment out steps 1-2 above)
-- CREATE INDEX IF NOT EXISTS products_name_fts_idx ON products USING GIN (to_tsvector('english', name));

-- 3. Optional: If you want to search across multiple fields (name, description, content, etc.)
-- You can create a combined tsvector column:
-- ALTER TABLE products ADD COLUMN IF NOT EXISTS search_vector tsvector
--     GENERATED ALWAYS AS (
--         setweight(to_tsvector('english', COALESCE(name, '')), 'A') ||
--         setweight(to_tsvector('english', COALESCE(description, '')), 'B') ||
--         setweight(to_tsvector('english', COALESCE(content, '')), 'C')
--     ) STORED;
-- 
-- CREATE INDEX IF NOT EXISTS products_search_vector_idx ON products USING GIN (search_vector);

-- Note: After running this migration, the full-text search in your Go code will work efficiently.
-- The 'plfts' operator in PostgREST uses plainto_tsquery which is more user-friendly for simple queries.
