DO $$
DECLARE
    password_hash bytea := '\x2432612431322449754751625575374c4b714b6d693435796632345575673863716d4c313057476768485a7149423739415969544239573952566b4f';
BEGIN
    INSERT INTO users (id, created_at, name, email, username, password_hash, activated, version)
    VALUES
        (uuid_generate_v4(), NOW(), 'John Doe', 'john.doe@example.com', 'johndoe', password_hash, true, 1),
        (uuid_generate_v4(), NOW(), 'Jane Smith', 'jane.smith@example.com', 'janesmith', password_hash, true, 1),
        (uuid_generate_v4(), NOW(), 'Alice Johnson', 'alice.johnson@example.com', 'alicejohnson', password_hash, true, 1),
        (uuid_generate_v4(), NOW(), 'Bob Brown', 'bob.brown@example.com', 'bobbrown', password_hash, true, 1),
        (uuid_generate_v4(), NOW(), 'Charlie Davis', 'charlie.davis@example.com', 'charliedavis', password_hash, true, 1),
        (uuid_generate_v4(), NOW(), 'Diana Evans', 'diana.evans@example.com', 'dianaevans', password_hash, true, 1),
        (uuid_generate_v4(), NOW(), 'Ethan Foster', 'ethan.foster@example.com', 'ethanfoster', password_hash, true, 1),
        (uuid_generate_v4(), NOW(), 'Fiona Green', 'fiona.green@example.com', 'fionagreen', password_hash, true, 1),
        (uuid_generate_v4(), NOW(), 'George Harris', 'george.harris@example.com', 'georgeharris', password_hash, true, 1),
        (uuid_generate_v4(), NOW(), 'Hannah White', 'hannah.white@example.com', 'hannahwhite', password_hash, true, 1);
END $$;
