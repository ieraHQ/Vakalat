-- Indexes for performance
CREATE INDEX idx_matters_client_id ON matters(client_id);
CREATE INDEX idx_matters_court_id ON matters(court_id);
CREATE INDEX idx_hearings_matter_id ON hearings(matter_id);
CREATE INDEX idx_documents_matter_id ON documents(matter_id);
CREATE INDEX idx_tasks_matter_id ON tasks(matter_id);
CREATE INDEX idx_invoices_matter_id ON invoices(matter_id);
CREATE INDEX idx_search_index_entity ON search_index(entity_type, entity_id);
CREATE INDEX idx_embeddings_document_id ON embeddings(document_id);
CREATE INDEX idx_embeddings_embedding ON embeddings USING ivfflat(embedding vector_cosine_ops);