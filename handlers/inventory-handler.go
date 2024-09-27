package handlers

import (
	"rizaldyaristyo-fiber-boiler/database"
	"rizaldyaristyo-fiber-boiler/models"

	"github.com/gofiber/fiber/v2"
)


func GetInventory(c *fiber.Ctx) error {
	manufacturingModels := []models.GetInventory{}
	inventoryRows, err := database.DB.Query(
		`
		-- omg my brain melted thinking this query 🫠, i admit it, im AI assisted with this one
		-- for this one query ill be documenting BETTER it so i won't lost in the sauce
		SELECT
			COALESCE(d.ManufacturingDemandQuantity, 0) AS ManufacturingDemandQuantity,
			COALESCE(d.SalesDemandQuantity, 0) AS SalesDemandQuantity,
			COALESCE(d.ManufacturingDemandQuantity, 0) + COALESCE(d.SalesDemandQuantity, 0) AS TotalDemandQuantity,
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
		LEFT JOIN vendors v ON inv.vendor_id = v.vendor_id
		LEFT JOIN (
			-- subquery for demand
			SELECT 
				inv.inventory_id,
				-- mo demand calculation (DISTINCT to avoid duplicates)
				SUM(DISTINCT mr.material_quantity_to_produce_product * mo.quantity) AS ManufacturingDemandQuantity,
				-- so demand calculation
				SUM(CASE WHEN s.sale_status = 'Pending' THEN s.quantity ELSE 0 END) AS SalesDemandQuantity
			FROM inventory inv
			LEFT JOIN manufacturing_recipes mr ON inv.inventory_id = mr.material_inventory_id
			LEFT JOIN manufacturing_orders mo ON mr.needed_to_produce_product_id = mo.product_id AND mo.status = 'Pending'
			LEFT JOIN sales s ON inv.inventory_id = s.item_id
			GROUP BY inv.inventory_id
		) d ON inv.inventory_id = d.inventory_id
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

func GetVendors(c *fiber.Ctx) error {
	var vendors models.GetVendors

	vendorRows, err := database.DB.Query("SELECT vendor_id, vendor_name, vendor_address, tax_id FROM vendors")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't find vendors", "data": nil})
	}

	for vendorRows.Next() {
		vendor := models.GetVendor{}

		vendorRows.Scan(
			&vendor.VendorID,
			&vendor.VendorName,
			&vendor.VendorAddress,
			&vendor.TaxID,
		)
		
		vendors = append(vendors, vendor)
	}

	defer vendorRows.Close()

	return c.JSON(vendors)
}

func ReplenishInventory(c *fiber.Ctx) error {
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

func EditInventory(c *fiber.Ctx) error {
	contentType := c.Get("Content-Type")
    if contentType == "application/x-www-form-urlencoded" || contentType == "multipart/form-data" {
		inventoryID := c.Params("inventory_id")

		edittedItemName := c.FormValue("edit-item-name")
		edittedVendorID := c.FormValue("edit-vendor-name") // the input called vendor name because it displays vendor name, tho the value is actually a vendor id
		edittedItemCode := c.FormValue("edit-item-code")
		edittedItemType := c.FormValue("edit-item-type")
		edittedSellable := c.FormValue("edit-sellable")
		edittedPurchasable := c.FormValue("edit-purchasable")
		edittedManufacturable := c.FormValue("edit-manufacturable")
		edittedPrice := c.FormValue("edit-price")
		edittedCurrency := c.FormValue("edit-currency")
		edittedQuantity := c.FormValue("edit-quantity")
		edittedQuantityWarning := c.FormValue("edit-quantity-warning")

		if (edittedItemName == ""){
			return c.Status(400).SendString("Missing edittedItemName")
		}
		if (edittedVendorID == ""){
			return c.Status(400).SendString("Missing edittedVendorID")
		}
		if (edittedItemCode == ""){
			return c.Status(400).SendString("Missing edittedItemCode")
		}
		if (edittedItemType == ""){
			return c.Status(400).SendString("Missing edittedItemType")
		}
		if (edittedSellable == ""){ edittedSellable = "0" } else { edittedSellable = "1"}
		if (edittedPurchasable == ""){ edittedPurchasable = "0" } else { edittedPurchasable = "1"}
		if (edittedManufacturable == ""){ edittedManufacturable = "0" } else { edittedManufacturable = "1"}
		if (edittedPrice == ""){
			return c.Status(400).SendString("Missing edittedPrice")
		}
		if (edittedCurrency == ""){
			return c.Status(400).SendString("Missing edittedCurrency")
		}
		if (edittedQuantity == ""){
			return c.Status(400).SendString("Missing edittedQuantity")
		}
		if (edittedQuantityWarning == ""){
			return c.Status(400).SendString("Missing edittedQuantityWarning")
		}
		
		// transaction
		tx, err := database.DB.Begin()
		if err != nil {
			return c.Status(500).SendString("Could not begin transaction")
		}

		_, err = tx.Exec(
			`
			UPDATE inventory
			SET 
				item_name = ?,
				vendor_id = ?,
				item_code = ?,
				item_type = ?,
				sellable = ?,
				purchaseable = ?,
				manufacturable = ?,
				price = ?,
				currency = ?,
				quantity = ?,
				minimum_stock_warning = ?
			WHERE inventory_id = ?
			`,
			edittedItemName,
			edittedVendorID,
			edittedItemCode,
			edittedItemType,
			edittedSellable,
			edittedPurchasable,
			edittedManufacturable,
			edittedPrice,
			edittedCurrency,
			edittedQuantity,
			edittedQuantityWarning,
			inventoryID,
		)
		if err != nil {
			tx.Rollback()
			return c.Status(500).SendString("Could not update inventory")
		}

		err = tx.Commit()
		if err != nil {
			return c.Status(500).SendString("Could not commit transaction")
		}

		return c.Status(200).SendString("OK")
	} else {
		return c.Status(400).SendString("Invalid content type")
	}
}