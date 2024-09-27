package handlers

import (
	"rizaldyaristyo-fiber-boiler/database"
	"rizaldyaristyo-fiber-boiler/models"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func OptimizedGetManufacturingOrder(c *fiber.Ctx) error {
	rows, err := database.DB.Query(
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
			mr.material_inventory_id AS MaterialID,
			mat.item_name AS MaterialName,
			mr.material_quantity_to_produce_product AS MaterialQuantityToProduceOne,
			mr.material_quantity_to_produce_product * mo.quantity AS TotalMaterialQuantityNeeded,
			mat.quantity AS MaterialCurrentQuantity
		FROM manufacturing_orders mo
		JOIN inventory prod 
			ON mo.product_id = prod.inventory_id
		LEFT JOIN manufacturing_recipes mr 
			ON mo.product_id = mr.needed_to_produce_product_id
		LEFT JOIN inventory mat 
			ON mr.material_inventory_id = mat.inventory_id
		WHERE mo.archived = false
		ORDER BY mo.order_id DESC;
		`,
	)

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	optimizedManufacturingOrderMap := make(map[int]*models.OptimizedGetManufacturingOrder)

	for rows.Next() {
		var optimizedManufacturingOrderModel models.OptimizedGetManufacturingOrder
		var optimizedManufacturingRecipe models.OptimizedManufacturingRecipe

		err := rows.Scan(
			&optimizedManufacturingOrderModel.OrderID,
			&optimizedManufacturingOrderModel.ManufactureOrderNumber,
			&optimizedManufacturingOrderModel.ProductID,
			&optimizedManufacturingOrderModel.ProductName,
			&optimizedManufacturingOrderModel.QuantityToManufacture,
			&optimizedManufacturingOrderModel.Status,
			&optimizedManufacturingOrderModel.CreatedAt,
			&optimizedManufacturingOrderModel.Archived,
			&optimizedManufacturingRecipe.MaterialID,
			&optimizedManufacturingRecipe.MaterialName,
			&optimizedManufacturingRecipe.MaterialQuantityToProduceOne,
			&optimizedManufacturingRecipe.TotalMaterialQuantityNeeded,
			&optimizedManufacturingRecipe.MaterialCurrentQuantity,
		)

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		// Check if the order already exists in the map
		if orderGroup, exists := optimizedManufacturingOrderMap[optimizedManufacturingOrderModel.OrderID]; exists {
			// Append the recipe to the existing order's Recipes field, but only if the material is valid
			if optimizedManufacturingRecipe.MaterialID.Valid {
				if orderGroup.Recipes == nil {
					orderGroup.Recipes = &[]models.OptimizedManufacturingRecipe{}
				}
				*orderGroup.Recipes = append(*orderGroup.Recipes, optimizedManufacturingRecipe)
			}
		} else {
			// Create a new order
			var recipes *[]models.OptimizedManufacturingRecipe
			if optimizedManufacturingRecipe.MaterialID.Valid {
				recipes = &[]models.OptimizedManufacturingRecipe{optimizedManufacturingRecipe} // Add the current recipe
			}

			newOrder := &models.OptimizedGetManufacturingOrder{
				OrderID:               optimizedManufacturingOrderModel.OrderID,
				ManufactureOrderNumber: optimizedManufacturingOrderModel.ManufactureOrderNumber,
				ProductID:             optimizedManufacturingOrderModel.ProductID,
				ProductName:           optimizedManufacturingOrderModel.ProductName,
				QuantityToManufacture: optimizedManufacturingOrderModel.QuantityToManufacture,
				Status:                optimizedManufacturingOrderModel.Status,
				CreatedAt:             optimizedManufacturingOrderModel.CreatedAt,
				Archived:              optimizedManufacturingOrderModel.Archived,
				Recipes:               recipes, // Set Recipes field to the pointer
			}
			optimizedManufacturingOrderMap[optimizedManufacturingOrderModel.OrderID] = newOrder
		}
	}


	optimizedManufacturingOrderModels := make([]models.OptimizedGetManufacturingOrder, 0, len(optimizedManufacturingOrderMap))
	for _, v := range optimizedManufacturingOrderMap {
		optimizedManufacturingOrderModels = append(optimizedManufacturingOrderModels, *v)
	}

	return c.JSON(optimizedManufacturingOrderModels)
}

func ApproveManufacturingOrder(c *fiber.Ctx) error {
	manufacturingOrderID := c.Params("manufacturing_order_id")

	// check if this needs material to produce
	var recipeExist bool
	var isPending bool
	err := database.DB.QueryRow(
		`SELECT EXISTS (SELECT 1 FROM manufacturing_recipes mr JOIN manufacturing_orders mo ON mr.needed_to_produce_product_id = mo.product_id WHERE mo.order_id = ?) AS need_recipe,
		EXISTS (SELECT 1 FROM manufacturing_orders WHERE status = "Pending" AND order_id = ?) AS is_pending;`,
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
			-- Deduce materials from inventory
			-- OPTIMIZED SO THE SERVER WON'T RUN A LOOP OF MULTIPLE UPDATES
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

		// Concatenate all parts to form the final query
		query := queryTop + queryMiddle + queryBottom

		// Prepare arguments for the CASE statement and WHERE clause
		args := []interface{}{}
		for _, materialSufficiency := range materialSufficiencies {
			args = append(args, materialSufficiency.MaterialID, materialSufficiency.NeededQuantity) // CASE arguments
		}
		for _, materialSufficiency := range materialSufficiencies {
			args = append(args, materialSufficiency.MaterialID) // WHERE clause arguments
		}

		// Execute the query
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

	// Updates MO to Completed
	_, err = tx.Exec(
		`
		UPDATE manufacturing_orders
		SET status = 'Completed'
		WHERE order_id = ?;
		`, manufacturingOrderID,
	)
	if err != nil {
		tx.Rollback()
		return c.Status(500).SendString(err.Error())
	}

	// Adds quantity to product manufactured
	_, err = tx.Exec(
		`
		UPDATE inventory i
		JOIN manufacturing_orders mo ON i.inventory_id = mo.product_id
		SET i.quantity = i.quantity + ?
		WHERE mo.order_id = ?;
		`, manufacturingQuantity, manufacturingOrderID,
	)

	if err != nil {
		tx.Rollback()
		return c.Status(500).SendString(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
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

func GetRecipes(c *fiber.Ctx) error {
	rows, err := database.DB.Query(`
		SELECT 
			inv_product.item_name AS ProductName,
			inv_product.item_code AS InventoryCode,  -- Product's inventory code
			inv_product.inventory_id AS ProductID,
			inv_material.inventory_id AS MaterialID,
			inv_material.item_name AS MaterialName,
			inv_material.item_code AS MaterialInventoryCode,  -- Material's inventory code
			mr.material_quantity_to_produce_product AS MaterialQuantityToProduceOne,
			inv_material.quantity AS MaterialCurrentQuantity
		FROM inventory inv_product
		LEFT JOIN manufacturing_recipes mr
			ON inv_product.inventory_id = mr.needed_to_produce_product_id
		LEFT JOIN inventory inv_material
			ON mr.material_inventory_id = inv_material.inventory_id
		WHERE inv_product.manufacturable = "1"
		ORDER BY inv_product.item_name, inv_material.item_name;

	`)

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	
	recipesMap := make(map[string]*models.GetProducts) // Map for fast lookups

	for rows.Next() {
		var getProductsModel models.GetProducts
		var getProductRecipeModel models.GetProductRecipe
		
		err := rows.Scan(
			&getProductsModel.ProductName,
			&getProductsModel.ProductInventoryCode,
			&getProductsModel.ProductID,
			&getProductRecipeModel.MaterialID,
			&getProductRecipeModel.MaterialName,
			&getProductRecipeModel.MaterialInventoryCode,
			&getProductRecipeModel.MaterialQuantityToProduceOne,
			&getProductRecipeModel.MaterialCurrentQuantity,
		)

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		// Check if product already exists in the map
		if productGroup, exists := recipesMap[strconv.Itoa(getProductsModel.ProductID)]; exists {
			// If MaterialID is valid, add the recipe
			if getProductRecipeModel.MaterialID.Valid {
				*productGroup.Recipes = append(*productGroup.Recipes, getProductRecipeModel)
			}
		} else {
			// Create a new product with its recipes
			var recipes *[]models.GetProductRecipe
			if getProductRecipeModel.MaterialID.Valid {
				recipes = &[]models.GetProductRecipe{getProductRecipeModel} // Initialize with the current recipe
			}

			newProduct := &models.GetProducts{
				ProductName:          getProductsModel.ProductName,
				ProductInventoryCode: getProductsModel.ProductInventoryCode,
				ProductID:            getProductsModel.ProductID,
				Recipes:              recipes,
			}

			// Use ProductID as the key for faster lookup
			recipesMap[strconv.Itoa(getProductsModel.ProductID)] = newProduct
		}
	}

	// Convert the map to a slice of GetProducts
	getProductsModels := make([]models.GetProducts, 0, len(recipesMap))
	for _, productGroup := range recipesMap {
		getProductsModels = append(getProductsModels, *productGroup)
	}

	return c.Status(200).JSON(getProductsModels)
}

func EditRecipe(c *fiber.Ctx) error {
    productInventoryID := c.Params("product_id")
    formData := c.Request().PostArgs()
    parsedForm := make(map[string][]string)

    formData.VisitAll(func(key, value []byte) {
        switch {
        case strings.HasPrefix(string(key), "recipe-material-current-quantity-"):
            parsedForm["recipe_material_current_quantity"] = append(parsedForm["recipe_material_current_quantity"], string(value))
        case strings.HasPrefix(string(key), "recipe-material-name-"):
            parsedForm["recipe_material_inventory_id"] = append(parsedForm["recipe_material_inventory_id"], string(value))
        case strings.HasPrefix(string(key), "recipe-material-quantity-to-produce-one-"):
            parsedForm["recipe_material_quantity_to_produce_one"] = append(parsedForm["recipe_material_quantity_to_produce_one"], string(value))
        }
    })

    // Begin transaction
    tx, err := database.DB.Begin()
    if err != nil {
        return c.Status(500).SendString(err.Error())
    }

    // Delete existing records
    _, err = tx.Exec(`
        DELETE FROM manufacturing_recipes WHERE needed_to_produce_product_id = ?;
    `, productInventoryID)

    if err != nil {
        tx.Rollback()
        return c.Status(500).SendString(err.Error())
    }

    // Construct bulk INSERT query
    queryTop := `
        INSERT INTO manufacturing_recipes 
        (needed_to_produce_product_id, material_inventory_id, material_quantity_to_produce_product) 
        VALUES 
    `
    queryValues := ""
    args := []interface{}{productInventoryID}

    for i := range parsedForm["recipe_material_inventory_id"] {
        if i > 0 {
            queryValues += ", "
        }
        queryValues += "(?, ?, ?)"
        args = append(args, parsedForm["recipe_material_inventory_id"][i], parsedForm["recipe_material_quantity_to_produce_one"][i])
    }

    query := queryTop + queryValues + ";"

    // Execute the query
    _, err = tx.Exec(query, args...)
    if err != nil {
        tx.Rollback()
        return c.Status(500).SendString(err.Error())
    }

    // Commit the transaction
    err = tx.Commit()
    if err != nil {
        return c.Status(500).SendString(err.Error())
    }

    return c.Status(200).SendString("OK")
}


func GetMaterials(c *fiber.Ctx) error {
	rows, err := database.DB.Query(`SELECT inventory_id, item_name, item_code FROM inventory WHERE item_type != "Product";`)

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	
	GetMaterialsModels := []models.GetMaterials{}

	for rows.Next() {
		var getMaterialsModel models.GetMaterials

		err := rows.Scan(
			&getMaterialsModel.InventoryID,
			&getMaterialsModel.ItemName,
			&getMaterialsModel.ItemCode,
		)

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		GetMaterialsModels = append(GetMaterialsModels, getMaterialsModel)
	}

	return c.Status(200).JSON(GetMaterialsModels)
}