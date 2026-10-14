package httpapi

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"uuid"

	"storeapi/src/internal/application"
	"storeapi/src/internal/domain"
)

type Handler struct {
	clients   *application.ClientService
	products  *application.ProductService
	suppliers *application.SupplierService
	images    *application.ImageService
}

func NewHandler(clients *application.ClientService, products *application.ProductService, suppliers *application.SupplierService, images *application.ImageService) *Handler {
	return &Handler{clients: clients, products: products, suppliers: suppliers, images: images}
}

// CreateClient godoc
// @Summary Create a client
// @Description Creates a client together with its address.
// @Resource clients
// @Accept json
// @Produce json
// @Param request body ClientCreateRequest true "Client to create"
// @Success 201 {object} ClientResponse "Client created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @OperationId createClient
// @Router /api/v1/clients [post]
func (h *Handler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var in ClientCreateRequest
	if err := decodeJSON(w, r, &in); err != nil {
		writeError(w, 400, err.Error())
		return
	}
	client, err := mapClientRequest(in)
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}
	client, err = h.clients.Create(r.Context(), client)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, mapClient(client))
}

// DeleteClient godoc
// @Summary Delete a client
// @Description Deletes a client by its UUID.
// @Resource clients
// @Param id path string true "Client UUID"
// @Success 204 "Client deleted successfully"
// @Failure 400 {object} ErrorResponse "Invalid UUID"
// @Failure 404 {object} ErrorResponse "Client not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @OperationId deleteClient
// @Router /api/v1/clients/{id} [delete]
func (h *Handler) DeleteClient(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}
	if err = h.clients.Delete(r.Context(), id); err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListClients godoc
// @Summary List or search clients
// @Description Returns all clients, optionally paginated, or searches clients by first and last name.
// @Resource clients
// @Produce json
// @Param first_name query string false "Client first name"
// @Param last_name query string false "Client last name"
// @Param limit query int false "Maximum number of clients to return"
// @Param offset query int false "Number of clients to skip"
// @Success 200 {array} ClientResponse "Clients"
// @Failure 400 {object} ErrorResponse "Invalid query parameters"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @OperationId listClients
// @Router /api/v1/clients [get]
func (h *Handler) ListClients(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var items []domain.Client
	var err error
	if q.Get("first_name") != "" || q.Get("last_name") != "" {
		items, err = h.clients.FindByName(r.Context(), q.Get("first_name"), q.Get("last_name"))
	} else {
		limit, offset, parseErr := pagination(q)
		if parseErr != nil {
			writeError(w, 400, parseErr.Error())
			return
		}
		items, err = h.clients.List(r.Context(), limit, offset)
	}
	if err != nil {
		handleError(w, err)
		return
	}
	out := make([]ClientResponse, 0, len(items))
	for _, v := range items {
		out = append(out, mapClient(v))
	}
	writeJSON(w, 200, out)
}

// ChangeClientAddress godoc
// @Summary Change client address
// @Description Replaces the current address of a client.
// @Resource clients
// @Accept json
// @Produce json
// @Param id path string true "Client UUID"
// @Param request body AddressDTO true "New client address"
// @Success 200 {object} ClientResponse "Updated client"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 404 {object} ErrorResponse "Client not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @OperationId changeClientAddress
// @Router /api/v1/clients/{id}/address [patch]
func (h *Handler) ChangeClientAddress(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}
	var in AddressDTO
	if err = decodeJSON(w, r, &in); err != nil {
		writeError(w, 400, err.Error())
		return
	}
	v, err := h.clients.ChangeAddress(r.Context(), id, mapAddress(in))
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, 200, mapClient(v))
}

// CreateProduct godoc
// @Summary Create a product
// @Description Creates a new appliance product.
// @Resource products
// @Accept json
// @Produce json
// @Param request body ProductCreateRequest true "Product to create"
// @Success 201 {object} ProductResponse "Product created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 404 {object} ErrorResponse "Supplier not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @OperationId createProduct
// @Router /api/v1/products [post]
func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var in ProductCreateRequest
	if err := decodeJSON(w, r, &in); err != nil {
		writeError(w, 400, err.Error())
		return
	}
	p, err := mapProductRequest(in)
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}
	p, err = h.products.Create(r.Context(), p)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, 201, mapProduct(p))
}

// AddProductPhoto godoc
// @Summary Add an existing photo to a product
// @Description Adds a photo to a product from the existing pool of photos
// @Resource products
// @Accept json
// @Produce json
// @Param id path string true "Product UUID"
// @Param request body AddProductPhotoRequest true "Add photo request"
// @Success 200 {object} ProductResponse "Updated product photo"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 404 {object} ErrorResponse "Product not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @OperationId addProductPhoto
// @Router /api/v1/products/{id}/image [put]
func (h *Handler) AddProductPhoto(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}
	var in AddProductPhotoRequest
	if err = decodeJSON(w, r, &in); err != nil {
		writeError(w, 400, err.Error())
		return
	}
	image, err := mapAddPhotoRequest(in)
	if err != nil {
		handleError(w, err)
		return
	}
	p, err := h.products.AddProductPhoto(r.Context(), id, image.ID)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, 200, mapProduct(p))
}

// DecreaseProductStock godoc
// @Summary Decrease product stock
// @Description Decreases the available stock of a product by the specified quantity.
// @Resource products
// @Accept json
// @Produce json
// @Param id path string true "Product UUID"
// @Param request body DecreaseStockRequest true "Stock decrease request"
// @Success 200 {object} ProductResponse "Updated product"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 404 {object} ErrorResponse "Product not found"
// @Failure 409 {object} ErrorResponse "Insufficient stock"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @OperationId decreaseProductStock
// @Router /api/v1/products/{id}/stock [patch]
func (h *Handler) DecreaseProductStock(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}
	var in DecreaseStockRequest
	if err = decodeJSON(w, r, &in); err != nil {
		writeError(w, 400, err.Error())
		return
	}
	p, err := h.products.DecreaseStock(r.Context(), id, in.Quantity)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, 200, mapProduct(p))
}

// GetProduct godoc
// @Summary Get product by ID
// @Description Returns a product by its UUID.
// @Resource products
// @Produce json
// @Param id path string true "Product UUID"
// @Success 200 {object} ProductResponse "Product"
// @Failure 400 {object} ErrorResponse "Invalid UUID"
// @Failure 404 {object} ErrorResponse "Product not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @OperationId getProduct
// @Router /api/v1/products/{id} [get]
func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}
	p, err := h.products.GetByID(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, 200, mapProduct(p))
}

// ListAvailableProducts godoc
// @Summary Get all available products
// @Description Returns all products whose available stock is greater than zero.
// @Resource products
// @Produce json
// @Success 200 {array} ProductResponse "Available products"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @OperationId listAvailableProducts
// @Router /api/v1/products/available [get]
func (h *Handler) ListAvailableProducts(w http.ResponseWriter, r *http.Request) {
	items, err := h.products.ListAvailable(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}
	out := make([]ProductResponse, 0, len(items))
	for _, v := range items {
		out = append(out, mapProduct(v))
	}
	writeJSON(w, 200, out)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Deletes a product by its UUID.
// @Resource products
// @Param id path string true "Product UUID"
// @Success 204 "Product deleted successfully"
// @Failure 400 {object} ErrorResponse "Invalid UUID"
// @Failure 404 {object} ErrorResponse "Product not found"
// @Failure 409 {object} ErrorResponse "Delete conflicts with related data"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @OperationId deleteProduct
// @Router /api/v1/products/{id} [delete]
func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}
	if err = h.products.Delete(r.Context(), id); err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(204)
}

// CreateSupplier godoc
// @Summary Create a supplier
// @Description Creates a supplier together with its address.
// @Resource suppliers
// @Accept json
// @Produce json
// @Param request body SupplierCreateRequest true "Supplier to create"
// @Success 201 {object} SupplierResponse "Supplier created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @OperationId createSupplier
// @Router /api/v1/suppliers [post]
func (h *Handler) CreateSupplier(w http.ResponseWriter, r *http.Request) {
	var in SupplierCreateRequest
	if err := decodeJSON(w, r, &in); err != nil {
		writeError(w, 400, err.Error())
		return
	}
	v, err := h.suppliers.Create(r.Context(), mapSupplierRequest(in))
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, 201, mapSupplier(v))
}

// ChangeSupplierAddress godoc
// @Summary Change supplier address
// @Description Replaces the current address of a supplier.
// @Resource suppliers
// @Accept json
// @Produce json
// @Param id path string true "Supplier UUID"
// @Param request body AddressDTO true "New supplier address"
// @Success 200 {object} SupplierResponse "Updated supplier"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 404 {object} ErrorResponse "Supplier not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @OperationId changeSupplierAddress
// @Router /api/v1/suppliers/{id}/address [patch]
func (h *Handler) ChangeSupplierAddress(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}
	var in AddressDTO
	if err = decodeJSON(w, r, &in); err != nil {
		writeError(w, 400, err.Error())
		return
	}
	v, err := h.suppliers.ChangeAddress(r.Context(), id, mapAddress(in))
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, 200, mapSupplier(v))
}

// DeleteSupplier godoc
// @Summary Delete a supplier
// @Description Deletes a supplier by its UUID.
// @Resource suppliers
// @Param id path string true "Supplier UUID"
// @Success 204 "Supplier deleted successfully"
// @Failure 400 {object} ErrorResponse "Invalid UUID"
// @Failure 404 {object} ErrorResponse "Supplier not found"
// @Failure 409 {object} ErrorResponse "Supplier cannot be deleted because related products exist"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @OperationId deleteSupplier
// @Router /api/v1/suppliers/{id} [delete]
func (h *Handler) DeleteSupplier(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}
	if err = h.suppliers.Delete(r.Context(), id); err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(204)
}

// ListSuppliers godoc
// @Summary Get all suppliers
// @Description Returns all suppliers.
// @Resource suppliers
// @Produce json
// @Success 200 {array} SupplierResponse "Suppliers"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @OperationId listSuppliers
// @Router /api/v1/suppliers [get]
func (h *Handler) ListSuppliers(w http.ResponseWriter, r *http.Request) {
	items, err := h.suppliers.List(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}
	out := make([]SupplierResponse, 0, len(items))
	for _, v := range items {
		out = append(out, mapSupplier(v))
	}
	writeJSON(w, 200, out)
}

// GetSupplier godoc
// @Summary Get supplier by ID
// @Description Returns a supplier by its UUID.
// @Resource suppliers
// @Produce json
// @Param id path string true "Supplier UUID"
// @Success 200 {object} SupplierResponse "Supplier"
// @Failure 400 {object} ErrorResponse "Invalid UUID"
// @Failure 404 {object} ErrorResponse "Supplier not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @OperationId getSupplier
// @Router /api/v1/suppliers/{id} [get]
func (h *Handler) GetSupplier(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}
	v, err := h.suppliers.GetByID(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, 200, mapSupplier(v))
}

// CreateImage godoc
// @Summary Add an image
// @Description Adds an image. The image field must contain the image bytes encoded as a Base64 string.
// @Resource images
// @Accept json
// @Produce json
// @Param request body ImageCreateRequest true "Product image"
// @Success 201 {object} ImageResponse "Created image identifiers"
// @Failure 400 {object} ErrorResponse "Invalid request or invalid Base64 data"
// @Failure 404 {object} ErrorResponse "Product not found"
// @Failure 409 {object} ErrorResponse "Product already has an image"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @OperationId createImage
// @Router /api/v1/images [post]
func (h *Handler) CreateImage(w http.ResponseWriter, r *http.Request) {
	var in ImageCreateRequest
	if err := decodeJSON(w, r, &in); err != nil {
		writeError(w, 400, err.Error())
		return
	}
	data, err := decodeImage(in.Image)
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}
	img, err := h.images.Create(r.Context(), domain.Image{Data: data})
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, ImageResponse{
		ID: img.ID.String(),
	})
}

// ReplaceImage godoc
// @Summary Replace an image
// @Description Replaces the binary data stored for an existing image. The image field must contain the image bytes encoded as a Base64 string.
// @Resource images
// @Accept json
// @Produce json
// @Param id path string true "Image UUID"
// @Param request body ImageReplaceRequest true "New image"
// @Success 200 {object} ImageResponse "Updated image identifiers"
// @Failure 400 {object} ErrorResponse "Invalid request or invalid Base64 data"
// @Failure 404 {object} ErrorResponse "Image not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @OperationId replaceImage
// @Router /api/v1/images/{id} [put]
func (h *Handler) ReplaceImage(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}
	var in ImageReplaceRequest
	if err = decodeJSON(w, r, &in); err != nil {
		writeError(w, 400, err.Error())
		return
	}
	data, err := decodeImage(in.Image)
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}
	img, err := h.images.Replace(r.Context(), id, data)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, ImageResponse{
		ID: img.ID.String(),
		// ProductID: img.ProductID.String(),
	})
}

// ListAvailableImages godoc
// @Summary Get all available images
// @Description Returns the entire pull of images
// @Resource images
// @Produce json
// @Success 200 {array} ProductResponse "Available images"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @OperationId listAvailableImages
// @Router /api/v1/images/available [get]
func (h *Handler) ListAvailableImages(w http.ResponseWriter, r *http.Request) {
	items, err := h.images.ListAvailable(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}
	out := make([]ImageResponse, 0, len(items))
	for _, v := range items {
		out = append(out, mapImage(v))
	}
	writeJSON(w, 200, out)
}

// DeleteImage godoc
// @Summary Delete an image
// @Description Deletes an image by its UUID.
// @Resource images
// @Param id path string true "Image UUID"
// @Success 204 "Image deleted successfully"
// @Failure 400 {object} ErrorResponse "Invalid UUID"
// @Failure 404 {object} ErrorResponse "Image not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @OperationId deleteImage
// @Router /api/v1/images/{id} [delete]
func (h *Handler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}
	if err = h.images.Delete(r.Context(), id); err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(204)
}

// GetProductImage godoc
// @Summary Get image of a product
// @Description Returns the binary image associated with a product.
// @Resource images
// @Produce application/octet-stream
// @Param product_id path string true "Product UUID"
// @Success 200 "Image bytes"
// @Failure 400 {object} ErrorResponse "Invalid UUID"
// @Failure 404 {object} ErrorResponse "Product image not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @OperationId getProductImage
// @Router /api/v1/products/{product_id}/image [get]
func (h *Handler) GetProductImage(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("product_id"))
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}
	img, err := h.images.GetByProductID(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}
	writeImage(w, img.ID, img.Data)
}

// GetImage godoc
// @Summary Get image by ID
// @Description Returns the binary image associated with an image UUID.
// @Resource images
// @Produce application/octet-stream
// @Param id path string true "Image UUID"
// @Success 200 "Image bytes"
// @Failure 400 {object} ErrorResponse "Invalid UUID"
// @Failure 404 {object} ErrorResponse "Image not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @OperationId getImage
// @Router /api/v1/images/{id} [get]
func (h *Handler) GetImage(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}
	img, err := h.images.GetByID(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}
	writeImage(w, img.ID, img.Data)
}

func writeImage(w http.ResponseWriter, id uuid.UUID, data []byte) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+id.String()+".bin\"")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(200)
	_, _ = w.Write(data)
}

func pagination(q url.Values) (*int, *int, error) {
	parse := func(key string) (*int, error) {
		raw := q.Get(key)
		if raw == "" {
			return nil, nil
		}
		v, err := strconv.Atoi(raw)
		if err != nil {
			return nil, errors.New(key + " must be an integer")
		}
		return &v, nil
	}
	limit, err := parse("limit")
	if err != nil {
		return nil, nil, err
	}
	offset, err := parse("offset")
	if err != nil {
		return nil, nil, err
	}
	return limit, offset, nil
}
