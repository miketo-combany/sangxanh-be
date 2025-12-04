# Full-Text Search Implementation

## Overview

This document explains the full-text search implementation for the product search functionality using PostgreSQL's built-in full-text search capabilities.

## What Changed

### Before (Simple LIKE Query)
```go
if filter.Name != "" {
    encoded := url.QueryEscape("%" + filter.Name + "%")
    query = query.Like("name", encoded)
}
```

### After (Full-Text Search)
```go
if filter.Name != "" {
    // Use full-text search with plainto_tsquery for better search results
    // This searches across the name field using PostgreSQL's FTS
    query = query.Filter("name", "plfts", filter.Name)
}
```

## Benefits of Full-Text Search

1. **Better Performance**: FTS uses GIN indexes which are much faster than LIKE queries on large datasets
2. **Language Awareness**: Handles stemming (searching "running" will match "run", "runs", "ran")
3. **Ranking**: Results can be ranked by relevance
4. **Word Boundaries**: Searches for whole words, not just substring matches
5. **Stopwords**: Automatically filters out common words like "the", "a", "an"
6. **Accent Insensitive**: Can handle accented characters properly with the right configuration

## Database Setup

### Step 1: Run the Migration

Execute the SQL migration file to add full-text search support:

```bash
# Using Supabase CLI
supabase db push migrations/add_fts_to_products.sql

# Or manually via SQL editor in Supabase dashboard
# Copy and paste the contents of migrations/add_fts_to_products.sql
```

### Step 2: Verify the Index

```sql
-- Check if the index was created
SELECT indexname, indexdef 
FROM pg_indexes 
WHERE tablename = 'products' AND indexname LIKE '%search%';
```

## PostgREST Full-Text Search Operators

The implementation uses the `plfts` operator which stands for "plain language full-text search". Here are the available operators:

- **`plfts`** (plainto_tsquery): Plain language search - converts user input to search query
  - Example: `"running shoes"` → searches for documents containing "run" and "shoe"
  
- **`fts`** (to_tsquery): More advanced but requires specific syntax
  - Example: `"running & shoes"` → must contain both words
  
- **`phfts`** (phraseto_tsquery): Phrase search - matches exact phrases
  - Example: `"running shoes"` → matches the exact phrase

- **`wfts`** (websearch_to_tsquery): Web-style search with quotes and operators
  - Example: `"running shoes" -cheap` → contains "running shoes" but not "cheap"

## Usage Examples

### Basic Search
```
GET /api/products?name=laptop
```
This will find all products with "laptop" in the name, including variations like "laptops".

### Multi-word Search
```
GET /api/products?name=gaming laptop
```
This will find products containing both "gaming" and "laptop".

## Advanced Configuration (Optional)

### Searching Multiple Fields

If you want to search across multiple fields (name, description, content), you can modify the migration to create a combined search vector:

```sql
ALTER TABLE products ADD COLUMN search_vector tsvector
    GENERATED ALWAYS AS (
        setweight(to_tsvector('english', COALESCE(name, '')), 'A') ||
        setweight(to_tsvector('english', COALESCE(description, '')), 'B') ||
        setweight(to_tsvector('english', COALESCE(content, '')), 'C')
    ) STORED;

CREATE INDEX products_search_vector_idx ON products USING GIN (search_vector);
```

Then update the Go code:
```go
if filter.Name != "" {
    query = query.Filter("search_vector", "plfts", filter.Name)
}
```

### Language Configuration

By default, the implementation uses English text search configuration. For other languages:

```sql
-- For Vietnamese
ALTER TABLE products ADD COLUMN name_search tsvector
    GENERATED ALWAYS AS (to_tsvector('simple', COALESCE(name, ''))) STORED;

-- Note: PostgreSQL doesn't have built-in Vietnamese support
-- Use 'simple' configuration for non-English languages without stemming
```

### Custom Ranking

To add relevance ranking to search results:

```go
// In the future, you could modify the query to include ranking
// This would require additional fields in the SELECT statement
query = query.Select("*, ts_rank(name_search, plainto_tsquery('english', $1)) as rank", filter.Name)
    .Filter("name_search", "@@", "plainto_tsquery('english', $1)")
    .OrderBy("rank", "desc")
```

## Performance Considerations

1. **Index Maintenance**: GIN indexes are updated automatically but can slow down INSERT/UPDATE operations slightly
2. **Storage**: The tsvector column adds storage overhead (typically 20-50% of text size)
3. **Memory**: GIN indexes require more memory than B-tree indexes
4. **Reindexing**: If you change language configuration, you need to rebuild the index

## Troubleshooting

### Search Returns No Results

1. Check if the migration was applied:
   ```sql
   \d products
   ```
   You should see the `name_search` column and index.

2. Verify the data is indexed:
   ```sql
   SELECT name, name_search FROM products LIMIT 5;
   ```

3. Test the search manually:
   ```sql
   SELECT name 
   FROM products 
   WHERE name_search @@ plainto_tsquery('english', 'your search term');
   ```

### Slow Queries

1. Check if the index is being used:
   ```sql
   EXPLAIN ANALYZE
   SELECT * FROM products 
   WHERE name_search @@ plainto_tsquery('english', 'laptop');
   ```
   Look for "Index Scan using products_name_search_idx"

2. Rebuild the index if needed:
   ```sql
   REINDEX INDEX products_name_search_idx;
   ```

## References

- [PostgreSQL Full-Text Search Documentation](https://www.postgresql.org/docs/current/textsearch.html)
- [PostgREST API Documentation](https://postgrest.org/en/stable/api.html#full-text-search)
- [Supabase Full-Text Search Guide](https://supabase.com/docs/guides/database/full-text-search)

## Next Steps

Consider implementing:
1. Search result highlighting
2. Search suggestions/autocomplete
3. Search analytics to track popular queries
4. Fuzzy matching for typos
5. Synonym support for better search coverage
