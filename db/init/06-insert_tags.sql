INSERT INTO tags (name, description, is_custom_tag) VALUES
    -- Political Leaning
    ('Progressive', 'Progressive', false),
    ('Liberal', 'Liberal', false),
    ('Moderate', 'Moderate', false),
    ('Conservative', 'Conservative', false),
    ('Libertarian', 'Libertarian', false),
    -- Professional Background
    ('Tech', 'Tech', false),
    ('Healthcare', 'Healthcare', false),
    ('Education', 'Education', false),
    ('Finance', 'Finance', false),
    ('Arts', 'Arts', false),
    ('Trade', 'Trade', false),
    ('Student', 'Student', false),
    -- Values
    ('Environmentalist', 'Environmentalist', false),
    ('Privacy Advocate', 'Privacy Advocate', false),
    ('Free Speech', 'Free Speech', false),
    ('Internationalist', 'Internationalist', false),
    ('Traditionalist', 'Traditionalist', false)
ON CONFLICT (name) DO NOTHING;  -- Using ON CONFLICT DO NOTHING makes this script safe to run multiple times should you need to