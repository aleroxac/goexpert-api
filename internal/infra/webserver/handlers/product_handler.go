package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aleroxac/goexpert-api/internal/dto"
	"github.com/aleroxac/goexpert-api/internal/entity"
	"github.com/aleroxac/goexpert-api/internal/infra/database"
	entityPkg "github.com/aleroxac/goexpert-api/pkg/entity"
	"github.com/go-chi/chi"
)

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: db,
	}
}

// CreateProduct godoc
//
//	@Summary		Create product
//	@Description	Create product
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			request	body	dto.ProductInput	true	"user credentials"
//	@Success		201
//	@Failure		500			{object}	dto.Error
//	@Failure		400			{object}	dto.Error
//	@Router			/products	[post]
//	@Security		ApiKeyAuth
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.ProductInput

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := dto.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	p, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := dto.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	err = h.ProductDB.Create(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := dto.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// GetProduct godoc
//
//	@Summary		Get product
//	@Description	Get product
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			id				path		string	true	"product ID"	Format(uuid)
//	@Success		200				{object}	entity.Product
//	@Failure		400				{object}	dto.Error
//	@Failure		404				{object}	dto.Error
//	@Router			/products/{id}	[get]
//	@Security		ApiKeyAuth
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		error := dto.Error{Message: "Invalid ID"}
		json.NewEncoder(w).Encode(error)
		return
	}

	_, err := entityPkg.ParseId(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := dto.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	product, err := h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		error := dto.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// GetProducts godoc
//
//	@Summary		List products
//	@Description	List products
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			page		query		string	false	"page number"
//	@Param			limit		query		string	false	"page limit"
//	@Success		200			{array}		entity.Product
//	@Failure		500			{object}	dto.Error
//	@Router			/products	[get]
//	@Security		ApiKeyAuth
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 0
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 0
	}
	sort := r.URL.Query().Get("sort")

	products, err := h.ProductDB.FindAll(pageInt, limitInt, sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := dto.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

// UpdateProduct godoc
//
//	@Summary		Update a product
//	@Description	Update a product
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			id				path		string				true	"product ID"	Format(uuid)
//	@Param			request			body		dto.ProductInput	true	"product request"
//	@Success		200				{object}	entity.Product
//	@Failure		400				{object}	dto.Error
//	@Failure		404				{object}	dto.Error
//	@Failure		500				{object}	dto.Error
//	@Router			/products/{id}	[put]
//	@Security		ApiKeyAuth
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var product entity.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		error := dto.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	product.ID, err = entityPkg.ParseId(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := dto.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	_, err = h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		error := dto.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	err = h.ProductDB.Update(&product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := dto.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// DeleteProduct godoc
//
//	@Summary		Delete a product
//	@Description	Delete a product
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			id				path		string		true	"product ID"	Format(uuid)
//	@Success		200
//	@Failure		400				{object}	dto.Error
//	@Failure		404				{object}	dto.Error
//	@Failure		500				{object}	dto.Error
//	@Router			/products/{id}	[delete]
//	@Security		ApiKeyAuth
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		error := dto.Error{Message: "Invalid ID"}
		json.NewEncoder(w).Encode(error)
		return
	}

	if _, err := entityPkg.ParseId(id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := dto.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	_, err := h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		error := dto.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	err = h.ProductDB.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := dto.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
