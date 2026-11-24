-- Universal holdout groups (RFC 0048, phase 4).
-- kind tells an experiment on a feature flag (0) from an internal experiment
-- that measures a holdout group (1). Every existing row is an experiment on a
-- feature flag.
ALTER TABLE `experiment` ADD COLUMN `kind` int NOT NULL DEFAULT 0, ALGORITHM=INSTANT;
