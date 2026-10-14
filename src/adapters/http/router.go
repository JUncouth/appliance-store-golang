package httpapi

import (
	"net/http"
)

func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("POST /api/v1/clients", h.CreateClient)
	mux.HandleFunc("GET /api/v1/clients", h.ListClients)
	mux.HandleFunc("DELETE /api/v1/clients/{id}", h.DeleteClient)
	mux.HandleFunc("PATCH /api/v1/clients/{id}/address", h.ChangeClientAddress)

	mux.HandleFunc("POST /api/v1/products", h.CreateProduct)
	mux.HandleFunc("PUT /api/v1/products/{id}/image", h.AddProductPhoto)
	mux.HandleFunc("GET /api/v1/products/available", h.ListAvailableProducts)
	mux.HandleFunc("GET /api/v1/products/{id}", h.GetProduct)
	mux.HandleFunc("PATCH /api/v1/products/{id}/stock", h.DecreaseProductStock)
	mux.HandleFunc("DELETE /api/v1/products/{id}", h.DeleteProduct)

	mux.HandleFunc("POST /api/v1/suppliers", h.CreateSupplier)
	mux.HandleFunc("GET /api/v1/suppliers", h.ListSuppliers)
	mux.HandleFunc("GET /api/v1/suppliers/{id}", h.GetSupplier)
	mux.HandleFunc("PATCH /api/v1/suppliers/{id}/address", h.ChangeSupplierAddress)
	mux.HandleFunc("DELETE /api/v1/suppliers/{id}", h.DeleteSupplier)

	mux.HandleFunc("POST /api/v1/images", h.CreateImage)
	mux.HandleFunc("PUT /api/v1/images/{id}", h.ReplaceImage)
	mux.HandleFunc("GET /api/v1/images/available", h.ListAvailableImages)
	mux.HandleFunc("DELETE /api/v1/images/{id}", h.DeleteImage)
	mux.HandleFunc("GET /api/v1/images/{id}", h.GetImage)
	mux.HandleFunc("GET /api/v1/products/{product_id}/image", h.GetProductImage)

	return logging(mux)
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { next.ServeHTTP(w, r) })
}
