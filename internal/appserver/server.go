package appserver

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

type Server struct {
	mux       *http.ServeMux
	provider  ChatProvider
	documents DocumentRepository
	storage   DocumentStorage
	retriever Retriever
	metrics   *RetrievalMetrics
}

func NewServer() *Server {
	return NewServerWithDependencies(MockProvider{}, NewMemoryDocumentRepository())
}

func NewServerWithDependencies(provider ChatProvider, documents DocumentRepository) *Server {
	return NewServerWithOptions(ServerOptions{
		Provider:  provider,
		Documents: documents,
	})
}

type ServerOptions struct {
	Provider  ChatProvider
	Documents DocumentRepository
	Storage   DocumentStorage
	Retriever Retriever
	Metrics   *RetrievalMetrics
}

func NewServerWithOptions(options ServerOptions) *Server {
	provider := options.Provider
	if provider == nil {
		provider = MockProvider{}
	}
	documents := options.Documents
	if documents == nil {
		documents = NewMemoryDocumentRepository()
	}
	metrics := options.Metrics
	if metrics == nil {
		metrics = NewRetrievalMetrics()
	}
	server := &Server{
		mux:       http.NewServeMux(),
		provider:  provider,
		documents: documents,
		storage:   options.Storage,
		retriever: options.Retriever,
		metrics:   metrics,
	}
	server.mux.HandleFunc("/healthz", server.healthz)
	server.mux.HandleFunc("/api/ops/retrieval-metrics", server.retrievalMetrics)
	server.mux.HandleFunc("/api/chat", server.chat)
	server.mux.HandleFunc("/api/documents", server.documentsRoot)
	server.mux.HandleFunc("/api/documents/", server.documentByID)
	return server
}

func (s *Server) retrievalMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	writeJSON(w, http.StatusOK, s.metrics.Snapshot())
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) healthz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func (s *Server) chat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var request ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}
	if strings.TrimSpace(request.Message) == "" {
		writeError(w, http.StatusBadRequest, "message is required")
		return
	}

	providerName := request.Provider
	if providerName == "" {
		providerName = "mock"
	}
	var retrieved RetrievalResult
	if request.RAG.Enabled && s.retriever != nil {
		var err error
		retrieved, err = s.retriever.Retrieve(r.Context(), RetrievalRequest{
			Message:     request.Message,
			DocumentIDs: request.RAG.DocumentIDs,
		})
		if err != nil {
			writeError(w, http.StatusBadGateway, "retrieval failed")
			return
		}
	}

	providerResponse, err := s.provider.Chat(ProviderRequest{
		Message:        request.Message,
		Context:        retrieved.Context,
		ClusterContext: request.ClusterContext,
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, "provider failed")
		return
	}

	writeJSON(w, http.StatusOK, ChatResponse{
		SessionID: firstNonEmpty(request.SessionID, "session-001"),
		MessageID: "msg-001",
		Answer:    providerResponse.Answer,
		Provider:  providerName,
		Citations: retrieved.Citations,
		Warnings:  []string{"Review AI generated content before use."},
	})
}

func (s *Server) documentsRoot(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, DocumentListResponse{Items: s.documents.List()})
	case http.MethodPost:
		s.uploadDocument(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *Server) uploadDocument(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "multipart form is required")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "file field is required")
		return
	}
	defer file.Close()

	if s.storage != nil {
		s.uploadDocumentToStorage(w, r, file, header)
		return
	}

	size, err := drainFile(file, header)
	if err != nil {
		writeError(w, http.StatusBadRequest, "could not read uploaded file")
		return
	}
	document := s.documents.Create(header.Filename, size, r.Header.Get("X-Forwarded-User"))
	writeJSON(w, http.StatusCreated, DocumentUploadResponse{
		ID:       document.ID,
		Filename: document.Filename,
		Status:   document.Status,
	})
}

func (s *Server) uploadDocumentToStorage(w http.ResponseWriter, r *http.Request, file multipart.File, header *multipart.FileHeader) {
	documentID, err := newUUID()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create document id")
		return
	}

	stored, err := s.storage.Store(r.Context(), documentID, header.Filename, file)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not store uploaded file")
		return
	}

	repository, ok := s.documents.(StoredDocumentRepository)
	if !ok {
		document := s.documents.Create(header.Filename, stored.SizeBytes, r.Header.Get("X-Forwarded-User"))
		writeJSON(w, http.StatusCreated, DocumentUploadResponse{
			ID:       document.ID,
			Filename: document.Filename,
			Status:   document.Status,
		})
		return
	}

	document, err := repository.CreateStored(r.Context(), CreateStoredDocumentInput{
		ID:         documentID,
		Filename:   header.Filename,
		SizeBytes:  stored.SizeBytes,
		ObjectURI:  stored.URI,
		UploadedBy: r.Header.Get("X-Forwarded-User"),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not persist document metadata")
		return
	}
	writeJSON(w, http.StatusCreated, DocumentUploadResponse{
		ID:       document.ID,
		Filename: document.Filename,
		Status:   document.Status,
	})
}

func (s *Server) documentByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/documents/")
	if id == "" || strings.Contains(id, "/") {
		writeError(w, http.StatusNotFound, "document not found")
		return
	}

	switch r.Method {
	case http.MethodGet:
		document, ok := s.documents.Get(id)
		if !ok {
			writeError(w, http.StatusNotFound, "document not found")
			return
		}
		writeJSON(w, http.StatusOK, document)
	case http.MethodDelete:
		document, ok := s.documents.MarkDeleting(id)
		if !ok {
			writeError(w, http.StatusNotFound, "document not found")
			return
		}
		writeJSON(w, http.StatusOK, DocumentDeleteResponse{ID: document.ID, Status: document.Status})
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func drainFile(file multipart.File, header *multipart.FileHeader) (int64, error) {
	if header.Size > 0 {
		_, err := io.Copy(io.Discard, file)
		return header.Size, err
	}
	return io.Copy(io.Discard, file)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
