ALTER TABLE orders DROP CONSTRAINT IF EXISTS fk_orders_document;
DROP TABLE IF EXISTS ocr_jobs;
DROP TABLE IF EXISTS document_tags;
DROP TABLE IF EXISTS embeddings;
DROP TABLE IF EXISTS document_versions;
DROP TABLE IF EXISTS documents;
