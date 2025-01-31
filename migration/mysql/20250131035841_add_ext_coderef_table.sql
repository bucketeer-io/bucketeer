ALTER TABLE code_reference
ADD COLUMN file_extension VARCHAR(32) NOT NULL DEFAULT '' COMMENT 'File extension of the code reference file'
AFTER file_path;
