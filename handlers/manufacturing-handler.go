package handlers

import (
	"rizaldyaristyo-fiber-boiler/database"
	"rizaldyaristyo-fiber-boiler/models"

	"github.com/gofiber/fiber/v2"
)

func GetManufacturingOrder(c *fiber.Ctx) error {
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

func ApproveManufacturingOrder(c *fiber.Ctx) error {
	manufacturingOrderID := c.Params("manufacturing_order_id")

	// check if this needs material to produce
	var recipeExist bool
	var isPending bool
	err := database.DB.QueryRow(
		`SELECT EXISTS (SELECT 1 FROM manufacturing_recipes WHERE needed_to_produce_product_id = ? ) AS need_recipe, EXISTS (SELECT 1 FROM manufacturing_orders where status = "Pending" and order_id = ?) AS is_pending;`,
		manufacturingOrderID, manufacturingOrderID,
	).Scan(&recipeExist, &isPending)
	
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	if !isPending {
		return c.Status(500).SendString("This manufacturing order is not pending")
	}
	
	if recipeExist {

		// check if material sufficient
		materialNeededRows, err := database.DB.Query(
			`
			SELECT 
				mr.material_inventory_id AS MaterialID,
				mat.item_name AS MaterialName,
				mat.quantity AS MaterialCurrentQuantity,
				mr.material_quantity_to_produce_product * mo.quantity AS TotalMaterialQuantityNeeded
			FROM manufacturing_orders mo
			JOIN inventory prod 
				ON mo.product_id = prod.inventory_id
			JOIN manufacturing_recipes mr 
				ON mo.product_id = mr.needed_to_produce_product_id
			JOIN inventory mat 
				ON mr.material_inventory_id = mat.inventory_id
			WHERE mo.status IN ('Pending', 'In Progress') AND mo.order_id = ?
			ORDER BY mo.created_at DESC;
			`, manufacturingOrderID,
		)

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		var materialSufficiencies models.MaterialSufficiencies
		
		for materialNeededRows.Next() {
			var materialSufficiency models.MaterialSufficiency
			err = materialNeededRows.Scan(
				&materialSufficiency.MaterialID,
				&materialSufficiency.MaterialName,
				&materialSufficiency.CurrentQuantity,
				&materialSufficiency.NeededQuantity,
			)
			if err != nil {
				return c.Status(500).SendString(err.Error())
			}

			if materialSufficiency.CurrentQuantity < materialSufficiency.NeededQuantity {
				return c.Status(500).SendString("Material not sufficient")
			}

			materialSufficiencies = append(materialSufficiencies, materialSufficiency)
		}

		// transaction
		tx, err := database.DB.Begin()

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		_, err = tx.Exec(
			`
			UPDATE manufacturing_orders
			SET status = 'In Progress'
			WHERE order_id = ?;
			`, manufacturingOrderID,
		)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		queryTop := `
			-- Deduce from inventory
			-- OPTIMIZED SO THE SERVER WONT RUN A LOOP OF MULTIPLE UPDATES

			UPDATE inventory
			SET quantity = quantity - CASE inventory_id
		`
		queryMiddle := ""
		for range materialSufficiencies {
			queryMiddle += " WHEN ? THEN ?"
		}
		queryBottom := `
				END
			WHERE inventory_id IN (
			`
		for i := range materialSufficiencies {
			if i > 0 {
				queryBottom += ", "
			}
			queryBottom += "?"
		}
		queryBottom += `);`
		query := queryTop + queryMiddle + queryBottom
		args := []interface{}{}
		for _, materialSufficiency := range materialSufficiencies {
			args = append(args, materialSufficiency.MaterialID, materialSufficiency.NeededQuantity)
		}
		for _, materialSufficiency := range materialSufficiencies {
			args = append(args, materialSufficiency.MaterialID)
		}
		_, err = tx.Exec(query, args...)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		// END OF OPTIMIZED SO THE SERVER WONT RUN A LOOP OF MULTIPLE UPDATES

		_, err = tx.Exec(query, args...)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		err = tx.Commit()
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

	} else {
		// transaction
		tx, err := database.DB.Begin()

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		_, err = tx.Exec(
			`
			UPDATE manufacturing_orders
			SET status = 'In Progress'
			WHERE order_id = ?;
			`, manufacturingOrderID,
		)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		err = tx.Commit()
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
	}
	
	return c.Status(200).SendString("OK")
}

func ReceiveManufacturingOrder(c *fiber.Ctx) error {
	manufacturingOrderID := c.Params("manufacturing_order_id")

	var manufacturingQuantity int
	var isInProgress bool

	// check how much being produced
	err := database.DB.QueryRow(`SELECT quantity, EXISTS( SELECT 1 FROM manufacturing_orders WHERE order_id = ? AND status = 'In Progress' ) FROM manufacturing_orders WHERE order_id = ?`, manufacturingOrderID, manufacturingOrderID).Scan(&manufacturingQuantity, &isInProgress)

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	if !isInProgress {
		return c.Status(500).SendString("Manufacturing not in progress")
	}

	// transaction
	tx, err := database.DB.Begin()

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	_, err = tx.Exec(
		`
		UPDATE manufacturing_orders
		SET status = 'Completed'
		WHERE order_id = ?;
		
		UPDATE inventory i
		JOIN manufacturing_orders mo ON i.inventory_id = mo.product_id
		SET i.quantity = i.quantity + ?
		WHERE mo.order_id = ?
		`, manufacturingOrderID, manufacturingQuantity, manufacturingOrderID,
	)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.Status(200).SendString("OK")
}

func CancelManufacturingOrder(c *fiber.Ctx) error {
	manufacturingOrderID := c.Params("manufacturing_order_id")

	// transaction
	tx, err := database.DB.Begin()

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	_, err = tx.Exec(
		`
		UPDATE manufacturing_orders
		SET status = 'Cancelled'
		WHERE order_id = ?;
		`, manufacturingOrderID,
	)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.Status(200).SendString("OK")
}