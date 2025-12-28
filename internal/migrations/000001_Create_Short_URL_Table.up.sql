-- Create extension for UUID generation if not exists
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create custom function to generate short codes
CREATE OR REPLACE FUNCTION generate_short_code(length INTEGER DEFAULT 6)
RETURNS TEXT AS $$
DECLARE
    chars TEXT := 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
    result TEXT := '';
    i INTEGER;
BEGIN
    FOR i IN 1..length LOOP
        result := result || substr(chars, floor(random() * length(chars) + 1)::INTEGER, 1);
    END LOOP;
    RETURN result;
END;
$$ LANGUAGE plpgsql;

-- Create function to ensure unique short codes
CREATE OR REPLACE FUNCTION generate_unique_short_code()
RETURNS TEXT AS $$
DECLARE
    new_code TEXT;
    max_attempts INTEGER := 100;
    attempt_count INTEGER := 0;
BEGIN
    LOOP
        attempt_count := attempt_count + 1;
        new_code := generate_short_code(6);
        
        -- Check if code already exists
        IF NOT EXISTS (SELECT 1 FROM short_urls WHERE short_code = new_code) THEN
            RETURN new_code;
        END IF;
        
        -- Prevent infinite loops
        IF attempt_count >= max_attempts THEN
            -- Try with longer code if too many collisions
            new_code := generate_short_code(8);
            IF NOT EXISTS (SELECT 1 FROM short_urls WHERE short_code = new_code) THEN
                RETURN new_code;
            END IF;
            
            -- Final fallback with UUID
            RETURN substr(replace(uuid_generate_v4()::TEXT, '-', ''), 1, 8);
        END IF;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- Create the short_urls table
CREATE TABLE short_urls (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    original_url TEXT NOT NULL,
    short_code VARCHAR(10) UNIQUE NOT NULL DEFAULT generate_unique_short_code(),
    title VARCHAR(255),
    description TEXT,
    click_count BIGINT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT valid_original_url CHECK (length(original_url) > 0),
    CONSTRAINT valid_short_code CHECK (length(short_code) >= 4),
    CONSTRAINT valid_expiry CHECK (expires_at IS NULL OR expires_at > created_at)
);

-- Create indexes for performance
CREATE INDEX idx_short_urls_short_code ON short_urls(short_code);
CREATE INDEX idx_short_urls_created_at ON short_urls(created_at);
CREATE INDEX idx_short_urls_expires_at ON short_urls(expires_at) WHERE expires_at IS NOT NULL;
CREATE INDEX idx_short_urls_active ON short_urls(is_active) WHERE is_active = true;
CREATE INDEX idx_short_urls_click_count ON short_urls(click_count);

-- Create function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for automatic updated_at updates
CREATE TRIGGER update_short_urls_updated_at
    BEFORE UPDATE ON short_urls
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Create function to increment click count
CREATE OR REPLACE FUNCTION increment_click_count(url_short_code VARCHAR(10))
RETURNS BOOLEAN AS $$
BEGIN
    UPDATE short_urls 
    SET click_count = click_count + 1,
        updated_at = CURRENT_TIMESTAMP
    WHERE short_code = url_short_code 
      AND is_active = true 
      AND (expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP);
    
    RETURN FOUND;
END;
$$ LANGUAGE plpgsql;

-- Create function to get URL by short code (with validation)
CREATE OR REPLACE FUNCTION get_url_by_short_code(url_short_code VARCHAR(10))
RETURNS TABLE (
    id UUID,
    original_url TEXT,
    title VARCHAR(255),
    is_expired BOOLEAN
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        su.id,
        su.original_url,
        su.title,
        CASE 
            WHEN su.expires_at IS NOT NULL AND su.expires_at <= CURRENT_TIMESTAMP THEN true
            ELSE false
        END as is_expired
    FROM short_urls su
    WHERE su.short_code = url_short_code 
      AND su.is_active = true;
END;
$$ LANGUAGE plpgsql;

-- Create function to clean up expired URLs
CREATE OR REPLACE FUNCTION cleanup_expired_urls()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM short_urls 
    WHERE expires_at IS NOT NULL 
      AND expires_at <= CURRENT_TIMESTAMP - INTERVAL '30 days';
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- Create a view for analytics (top URLs)
CREATE VIEW top_urls AS
SELECT 
    id,
    original_url,
    short_code,
    title,
    click_count,
    created_at,
    CASE 
        WHEN expires_at IS NOT NULL AND expires_at <= CURRENT_TIMESTAMP THEN true
        ELSE false
    END as is_expired
FROM short_urls
WHERE is_active = true
ORDER BY click_count DESC, created_at DESC;