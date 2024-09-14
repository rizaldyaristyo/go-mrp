package handlers

import (
	"rizaldyaristyo-fiber-boiler/database"
	"rizaldyaristyo-fiber-boiler/models"

	"github.com/gofiber/fiber/v2"
)


func GetInventory(c *fiber.Ctx) error {
	permissionInt := getNthDigit(c.Locals("permission_val").(int), 0)
	if permissionInt < 0 {
		return c.Status(401).SendString("Access Denied")
	}

	manufacturingModels := []models.GetInventory{}
	inventoryRows, err := database.DB.Query(
		`
		-- omg my brain melted thinking this query 🫠, i admit it, im assisted by AI with this one
		-- for this one query ill be documenting it so i won't lost in the sauce
		SELECT 
			COALESCE(SUM(mr.material_quantity_to_produce_product * mo.quantity), 0) AS ManufacturingDemandQuantity,
			COALESCE(SUM(s.quantity), 0) AS SalesDemandQuantity,
			COALESCE(SUM(mr.material_quantity_to_produce_product * mo.quantity), 0) +
			COALESCE(SUM(s.quantity), 0) AS TotalDemandQuantity,
			inv.inventory_id AS InventoryID,
			inv.item_name AS ItemName,
			inv.vendor_id AS VendorID,
			v.vendor_name AS VendorName,
			v.vendor_address AS VendorAddress,
			v.tax_id AS TaxID,
			inv.item_code AS ItemCode,
			inv.item_code_2 AS ItemCode2,
			inv.item_type AS ItemType,
			inv.sellable AS Sellable,
			inv.purchaseable AS Purchaseable,
			inv.manufacturable AS Manufacturable,
			inv.price AS Price,
			inv.price_2 AS Price2,
			inv.currency AS Currency,
			inv.quantity AS Quantity,
			inv.minimum_stock_warning AS MinimumStock,
			inv.last_updated AS LastUpdated,
			inv.archived AS Archived
		FROM inventory inv
		LEFT JOIN vendors v 
			ON inv.vendor_id = v.vendor_id
		LEFT JOIN manufacturing_recipes mr 
			ON inv.inventory_id = mr.material_inventory_id
		LEFT JOIN manufacturing_orders mo 
			ON mr.needed_to_produce_product_id = mo.product_id
			AND mo.status IN ('Pending', 'In Progress')
		LEFT JOIN sales s 
			ON inv.inventory_id = s.item_id
			AND s.sale_status IN ('Pending', 'Paid')
		GROUP BY 
			inv.inventory_id, 
			inv.item_name, 
			inv.vendor_id, 
			v.vendor_name, 
			v.vendor_address, 
			v.tax_id,
			inv.item_code, 
			inv.item_code_2, 
			inv.item_type, 
			inv.sellable, 
			inv.purchaseable, 
			inv.manufacturable,
			inv.price, 
			inv.price_2, 
			inv.currency, 
			inv.quantity, 
			inv.minimum_stock_warning, 
			inv.last_updated, 
			inv.archived
		ORDER BY TotalDemandQuantity DESC;
		`,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error 1", "message": "Couldn't find inventory", "data": nil})
	}

	for inventoryRows.Next() {

		inventory := models.GetInventory{}
		
		inventoryRows.Scan(
			&inventory.ManufacturingDemandQuantity,
			&inventory.SalesDemandQuantity,
			&inventory.TotalDemandQuantity,
			&inventory.InventoryID,
			&inventory.ItemName,
			&inventory.VendorID,
			&inventory.VendorName,
			&inventory.VendorAddress,
			&inventory.TaxID,
			&inventory.ItemCode,
			&inventory.ItemCode2,
			&inventory.ItemType,
			&inventory.Sellable,
			&inventory.Purchaseable,
			&inventory.Manufacturable,
			&inventory.Price,
			&inventory.Price2,
			&inventory.Currency,
			&inventory.Quantity,
			&inventory.MinimumStock,
			&inventory.LastUpdated,
			&inventory.Archived,
		)

		manufacturingModels = append(manufacturingModels, inventory)
	}
	
	defer inventoryRows.Close()
	
	return c.JSON(manufacturingModels)
}

func ReplenishInventory(c *fiber.Ctx) error {
	permissionInt := getNthDigit(c.Locals("permission_val").(int), 0)
	if permissionInt < 0 {
		return c.Status(401).SendString("Access Denied")
	}

	contentType := c.Get("Content-Type")
    if contentType == "application/x-www-form-urlencoded" || contentType == "multipart/form-data" {
		inventoryID := c.Params("inventory_id")

		replenishOrderNumber := c.FormValue("replenish-order-number")
		if replenishOrderNumber == "" {
			return c.Status(400).SendString("Missing order number")
		}

		quantity := c.FormValue("replenish-quantity")
		if quantity == "" {
			return c.Status(400).SendString("Missing quantity")
		}

		switch c.FormValue("replenish-manufacture-or-purchase-radio") {
		case "manufacture":
			tx, err := database.DB.Begin()
			if err != nil {
				return c.Status(500).SendString("Could not begin transaction")
			}

			// return c.Status(200).SendString("manufacture_order_number: " + replenishOrderNumber + " inventory_id: " + inventoryID + " quantity: " + quantity + " status: " + "Pending")

			_, err = tx.Exec("INSERT INTO manufacturing_orders (manufacture_order_number, product_id, quantity, status) VALUES (?, ?, ?, ?)", replenishOrderNumber, inventoryID, quantity, "Pending")
			
			if err != nil {
				tx.Rollback()
				return c.Status(500).SendString("Could not insert into manufacturing_orders" + err.Error())
			}

			err = tx.Commit()
			if err != nil {
				return c.Status(500).SendString("Could not commit transaction")
			}

			// return c.Render("inventory", fiber.Map{
			// 	"inventoryMessage": "Manufacture order created",
			// })

			return c.Redirect("/inventory?replenish-success=true&replenish-type=manufacture")
			
		case "purchase":

			// check item price
			var itemPrice float64
			err := database.DB.QueryRow("SELECT price FROM inventory WHERE inventory_id = ?", inventoryID).Scan(&itemPrice)

			if err != nil {
				return c.Status(500).SendString("Could not check item price")
			}

			tx, err := database.DB.Begin()
			if err != nil {
				return c.Status(500).SendString("Could not begin transaction")
			}

			_, err = tx.Exec("INSERT INTO purchases (purchase_order_number, item_id, quantity, purchase_price_per_unit, purchase_status) VALUES (?, ?, ?, ?, ?)", replenishOrderNumber, inventoryID, quantity, itemPrice, "Pending")

			if err != nil {
				tx.Rollback()
				return c.Status(500).SendString("Could not insert into purchase_orders")
			}

			err = tx.Commit()
			if err != nil {
				return c.Status(500).SendString("Could not commit transaction")
			}

			return c.Redirect("/inventory?replenish-success=true&replenish-type=purchase")
		default:
			return c.Status(400).SendString("Missing manufacture or purchase")
		}
	} else {
		return c.Status(400).SendString("Invalid content type")
	}
}