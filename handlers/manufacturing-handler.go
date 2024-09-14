package handlers

import (
	"rizaldyaristyo-fiber-boiler/database"
	"rizaldyaristyo-fiber-boiler/models"

	"github.com/gofiber/fiber/v2"
)

func GetManufacturingOrder(c *fiber.Ctx) error {
	permissionInt := getNthDigit(c.Locals("permission_val").(int), 0)
	if permissionInt < 0 {
		return c.Status(401).SendString("Access Denied")
	}

	manufacturingOrderModels := []models.GetManufacturingOrder{}

	manufacturingOrderRows, err := database.DB.Query(
		`SELECT
			mo.order_id AS OrderID,
			mo.manufacture_order_number AS ManufactureOrderNumber,
			mo.product_id AS ProductID,
			inv.item_name AS ProductName,
			mo.quantity AS QuantityToManufacture,
			mo.status AS Status,
			mo.created_at AS CreatedAt,
			mo.archived AS Archived
		FROM manufacturing_orders mo
		JOIN inventory inv 
			ON mo.product_id = inv.inventory_id
		WHERE mo.archived = false
		ORDER BY mo.created_at DESC;
		`,
	)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}


	for manufacturingOrderRows.Next() {
		var manufacturingOrderModel models.GetManufacturingOrder
		err = manufacturingOrderRows.Scan(
			&manufacturingOrderModel.OrderID,
			&manufacturingOrderModel.ManufactureOrderNumber,
			&manufacturingOrderModel.ProductID,
			&manufacturingOrderModel.ProductName,
			&manufacturingOrderModel.QuantityToManufacture,
			&manufacturingOrderModel.Status,
			&manufacturingOrderModel.CreatedAt,
			&manufacturingOrderModel.Archived,
		)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}


		manufacturingRecipeRows, err := database.DB.Query(
			`
			SELECT 
				mo.order_id AS OrderID,
				mo.manufacture_order_number AS ManufactureOrderNumber,
				mo.product_id AS ProductID,
				prod.item_name AS ProductName,
				mo.quantity AS QuantityToManufacture,
				mo.status AS Status,
				mo.created_at AS CreatedAt,
				mo.archived AS Archived,
				mr.material_quantity_to_produce_product AS MaterialQuantityToProduceOne,
				mr.material_inventory_id AS MaterialID,
				mat.item_name AS MaterialName,
				mr.material_quantity_to_produce_product * mo.quantity AS TotalMaterialQuantityNeeded,
				mat.quantity AS MaterialCurrentQuantity
			FROM manufacturing_orders mo
			JOIN inventory prod 
				ON mo.product_id = prod.inventory_id
			JOIN manufacturing_recipes mr 
				ON mo.product_id = mr.needed_to_produce_product_id
			JOIN inventory mat 
				ON mr.material_inventory_id = mat.inventory_id
			WHERE mo.status IN ('Pending', 'In Progress') 
			AND mo.manufacture_order_number = ?
			ORDER BY mo.created_at DESC;
			`, manufacturingOrderModel.ManufactureOrderNumber,
		)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		for manufacturingRecipeRows.Next() {
			var manufacturingRecipeModel models.ManufacturingRecipe

			err = manufacturingRecipeRows.Scan(
				&manufacturingRecipeModel.OrderID,
				&manufacturingRecipeModel.ManufactureOrderNumber,
				&manufacturingRecipeModel.ProductID,
				&manufacturingRecipeModel.ProductName,
				&manufacturingRecipeModel.QuantityToManufacture,
				&manufacturingRecipeModel.Status,
				&manufacturingRecipeModel.CreatedAt,
				&manufacturingRecipeModel.Archived,
				&manufacturingRecipeModel.MaterialQuantityToProduceOne,
				&manufacturingRecipeModel.MaterialID,
				&manufacturingRecipeModel.MaterialName,
				&manufacturingRecipeModel.TotalMaterialQuantityNeeded,
				&manufacturingRecipeModel.MaterialCurrentQuantity,
			)
			if err != nil {
				return c.Status(500).SendString(err.Error())
			}

			manufacturingOrderModel.Recipes = append(manufacturingOrderModel.Recipes, manufacturingRecipeModel)
		}

		manufacturingOrderModels = append(manufacturingOrderModels, manufacturingOrderModel)
	}

	return c.JSON(manufacturingOrderModels)
}