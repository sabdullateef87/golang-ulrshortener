-- Drop view
DROP VIEW IF EXISTS top_urls;

-- Drop functions
DROP FUNCTION IF EXISTS cleanup_expired_urls();
DROP FUNCTION IF EXISTS get_url_by_short_code(VARCHAR(10));
DROP FUNCTION IF EXISTS increment_click_count(VARCHAR(10));
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP FUNCTION IF EXISTS generate_unique_short_code();
DROP FUNCTION IF EXISTS generate_short_code(INTEGER);

-- Drop table (this will automatically drop indexes and triggers)
DROP TABLE IF EXISTS short_urls;

-- Note: We don't drop the uuid-ossp extension as it might be used by other parts of the application
-- If you want to drop it uncomment the line below:
-- DROP EXTENSION IF EXISTS "uuid-ossp";