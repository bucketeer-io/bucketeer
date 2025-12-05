-- ==========================================
-- UUID v4 Migration for api_key table
-- ==========================================

-- 1. Create uuid_v4() function
DELIMITER $$

DROP FUNCTION IF EXISTS uuid_v4$$

CREATE FUNCTION uuid_v4()
    RETURNS CHAR(36)
    NOT DETERMINISTIC
BEGIN
    DECLARE b BINARY(16);

    -- 16 random bytes
    SET b = RANDOM_BYTES(16);

    -- Set version to 4 (0100xxxx at byte 7)
    SET b = INSERT(
            b, 7, 1,
            CHAR((ASCII(SUBSTR(b, 7, 1)) & 0x0F) | 0x40)
            );

    -- Set variant to 10xxxxxx at byte 9
    SET b = INSERT(
            b, 9, 1,
            CHAR((ASCII(SUBSTR(b, 9, 1)) & 0x3F) | 0x80)
            );

    RETURN BIN_TO_UUID(b);
END$$

DELIMITER ;

-- 2. Update all api_key IDs to UUID v4
UPDATE api_key
SET id = uuid_v4();

-- 3. Cleanup (optional - remove function after migration)
DROP FUNCTION IF EXISTS uuid_v4;
