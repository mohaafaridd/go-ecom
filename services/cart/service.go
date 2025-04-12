package cart

import (
	"fmt"

	"mohaafaridd.dev/ecom/types"
)

func getCartItemsIds(items []types.CartItem) ([]int, error) {
	productIds := make([]int, len(items))

	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for the product %d", item.ProductId)
		}

		productIds[i] = item.ProductId
	}

	return productIds, nil
}

func (h *Handler) createOrder(products []types.Product, items []types.CartItem, userId int) (int, float64, error) {
	productMap := make(map[int]types.Product)

	for _, product := range products {
		productMap[product.ID] = product
	}

	// Check product availability
	if err := checkStock(items, productMap); err != nil {
		return 0, 0, err
	}
	// Calculate total price
	totalPrice := calculateTotalPrice(items, productMap)

	// Reduce product quantity
	for _, item := range items {
		product := productMap[item.ProductId]
		product.Quantity -= item.Quantity

		h.productStore.UpdateProductQuantity(product.ID, product.Quantity)
	}

	// create the order
	orderId, err := h.store.CreateOrder(types.Order{
		UserId:  userId,
		Total:   totalPrice,
		Status:  "pending",
		Address: "placeholder",
	})

	if err != nil {
		return 0, 0, err
	}

	// create order items
	for _, item := range items {
		h.store.CreateOrderItem(types.OrderItem{
			OrderId:   orderId,
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
			Price:     productMap[item.ProductId].Price,
		})
	}

	return orderId, totalPrice, nil
}

func checkStock(items []types.CartItem, productMap map[int]types.Product) error {
	if len(items) == 0 {
		return fmt.Errorf("Empty cart")
	}

	for _, item := range items {
		product, ok := productMap[item.ProductId]

		if !ok {
			return fmt.Errorf("product %d isn't available in store, please refresh your cart", item.ProductId)
		}

		if product.Quantity < item.Quantity {
			return fmt.Errorf("product %d isn't available in the quantity requested", item.ProductId)

		}
	}

	return nil
}

func calculateTotalPrice(items []types.CartItem, productMap map[int]types.Product) float64 {
	var total float64

	for _, item := range items {
		total += productMap[item.ProductId].Price * float64(item.Quantity)
	}

	return total
}
