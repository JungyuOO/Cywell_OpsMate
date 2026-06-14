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
}

func NewServer() *Server {
	return NewServerWithDependencies(MockProvider{}, NewMemoryDocumentRepository())
}

func NewServerWithDependencies(provider ChatProvider, documents DocumentRepository) *Server {
	server := &Server{
		mux:       http.NewServeMux(),
		provider:  provider,
		documents: documents,
	}
	server.mux.HandleFunc("/healthz", server.healthz)
	server.mux.HandleFunc("/api/chat", server.chat)
	server.mux.HandleFunc("/api/documents", server.documentsRoot)
	server.mux.HandleFunc("/api/documents/", server.documentByID)
	return server
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
	providerResponse, err := s.provider.Chat(ProviderRequest{
		Message:        request.Message,
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
