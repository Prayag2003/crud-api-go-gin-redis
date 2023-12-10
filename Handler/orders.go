package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	model "github.com/Prayag2003/crud-api-golang/Model"
	order "github.com/Prayag2003/crud-api-golang/Repository/Order"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Order struct {
	Repo *order.RedisRepo
}

var errNotExist error = errors.New("order does not exist")

func (o *Order) Create(c *gin.Context) {
	var body struct {
		CustomerID uuid.UUID    `json:"customer_id"`
		Items      []model.Item `json:"items"`
	}

	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()
	order := model.Order{
		OrderID:    rand.Uint64(),
		CustomerID: body.CustomerID,
		Items:      body.Items,
		CreatedAt:  &now,
	}

	err := o.Repo.Insert(c.Request.Context(), order)

	if err != nil {
		fmt.Println("Failed to insert data")
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(order)
	if err != nil {
		fmt.Println("failed to marshall:", err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
	}
	c.Writer.Write(res)
	c.Writer.WriteHeader(http.StatusCreated)
}

func (o *Order) List(c *gin.Context) {
	r := c.Request
	w := c.Writer
	cursorStr := r.URL.Query().Get("cursor")
	if cursorStr == "" {
		cursorStr = "0"
	}

	const decimal = 10
	const bitSize = 64
	cursor, err := strconv.ParseUint(cursorStr, decimal, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	const size = 50
	res, err := o.Repo.FindAll(r.Context(), order.FindAllPage{
		Offset: cursor,
		Size:   size,
	})
	if err != nil {
		fmt.Println("failed to find all:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var response struct {
		Items []model.Order `json:"items"`
		Next  uint64        `json:"next,omitempty"`
	}
	response.Items = res.Orders
	response.Next = res.Cursor

	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func (o *Order) OrderById(c *gin.Context) {
	w := c.Writer
	r := c.Request
	idParam := c.Param("id")

	const base = 10
	const bitSize = 64

	orderID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h, err := o.Repo.FindByID(r.Context(), orderID)
	if errors.Is(err, order.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find by id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(h); err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (o *Order) UpdateById(c *gin.Context) {
	var body struct {
		Status string `json:"status"`
	}
	w := c.Writer
	r := c.Request
	idParam := c.Param("id")

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	const base = 10
	const bitSize = 64

	orderID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	theOrder, err := o.Repo.FindByID(r.Context(), orderID)
	if errors.Is(err, order.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find by id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	const completedStatus = "completed"
	const shippedStatus = "shipped"
	now := time.Now().UTC()

	switch body.Status {
	case shippedStatus:
		if theOrder.ShippedAt != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		theOrder.ShippedAt = &now
	case completedStatus:
		if theOrder.CompletedAt != nil || theOrder.ShippedAt == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		theOrder.CompletedAt = &now
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = o.Repo.Update(r.Context(), theOrder)
	if err != nil {
		fmt.Println("failed to insert:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(theOrder); err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (o *Order) DeleteById(c *gin.Context) {
	w := c.Writer
	r := c.Request
	idParam := c.Param("id")

	const base = 10
	const bitSize = 64

	orderID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = o.Repo.DeleteByID(r.Context(), orderID)
	if errors.Is(err, order.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find by id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
