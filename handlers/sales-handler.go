package handlers

import (
	"log"
	"rizaldyaristyo-fiber-boiler/database"
	"rizaldyaristyo-fiber-boiler/models"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetSalesSensitive(c *fiber.Ctx) error {
	salesRows, err := database.DB.Query(
		`
		SELECT
			s.sales_order_number AS 'SalesOrderNumber',
			i.item_name AS 'ProductName',
			s.quantity AS 'Quantity',
			s.sent_quantity AS 'SentQuantity',
			i.price AS 'MfgPricePerUnit',
			s.sale_price_per_unit AS 'SalePriceperUnit',
			s.tax_percent AS 'TaxPercent',
			c.customer_id AS 'CustomerID',
			c.customer_name AS 'CustomerName',
			c.customer_bank_name AS 'CustomerBankName',
			c.customer_bank_account_number AS 'CustomerBankAccountNumber',
			c.customer_cc_number AS 'CustomerCCNumber',
			c.customer_address AS 'CustomerAddress',
			c.customer_email AS 'CustomerEmail',
			c.customer_phone AS 'CustomerPhone',
			c.customer_tax_id AS 'CustomerTaxID',
			s.sales_channel AS 'SalesChannel',
			s.payment_method AS 'PaymentMethod',
			s.payment_status AS 'PaymentStatus',
			s.delivery_status AS 'DeliveryStatus',
			s.canceled AS 'IsCanceled',
			s.created_at AS 'OrderDate',
			s.payment_date AS 'PaymentDate',
			s.delivery_date AS 'DeliveryDate'
		FROM sales s
		JOIN customers c ON s.customer_id = c.customer_id
		JOIN inventory i ON s.item_id = i.inventory_id
		WHERE i.sellable = 1
		ORDER BY s.created_at DESC;
		`)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't find sales", "data": nil})
	}

	salesModels := []models.GetSaleSensitive{}

	for salesRows.Next() {
		salesModel := models.GetSaleSensitive{}
		customerModel := models.GetCustomerSensitive{}

		salesRows.Scan(
			&salesModel.SalesOrderNumber,
			&salesModel.ProductName,
			&salesModel.Quantity,
			&salesModel.SentQuantity,
			&salesModel.MfgPricePerUnit,
			&salesModel.SalePriceperUnit,
			&salesModel.TaxPercent,
			&customerModel.CustomerID,
			&customerModel.CustomerName,
			&customerModel.CustomerBankName,
			&customerModel.CustomerBankAccountNumber,
			&customerModel.CustomerCCNumber,
			&customerModel.CustomerAddress,
			&customerModel.CustomerEmail,
			&customerModel.CustomerPhone,
			&customerModel.CustomerTaxID,
			&salesModel.SalesChannel,
			&salesModel.PaymentMethod,
			&salesModel.PaymentStatus,
			&salesModel.DeliveryStatus,
			&salesModel.IsCanceled,
			&salesModel.OrderDate,
			&salesModel.PaymentDate,
			&salesModel.DeliveryDate,
		)

		salesModel.Customer = customerModel
		salesModels = append(salesModels, salesModel)
	}

	defer salesRows.Close()

	return c.JSON(salesModels)
}





func OptimizedGetSalesSensitive(c *fiber.Ctx) error {
    salesRows, err := database.DB.Query(
        `
        SELECT
            s.sales_order_number AS 'SalesOrderNumber',
            i.inventory_id AS 'ProductID',
            i.item_name AS 'ProductName',
            s.quantity AS 'Quantity',
            s.sent_quantity AS 'SentQuantity',
            i.price AS 'MfgPricePerUnit',
            s.sale_price_per_unit AS 'SalePriceperUnit',
            s.tax_percent AS 'TaxPercent',
            c.customer_id AS 'CustomerID',
            c.customer_name AS 'CustomerName',
            c.customer_bank_name AS 'CustomerBankName',
            c.customer_bank_account_number AS 'CustomerBankAccountNumber',
            c.customer_cc_number AS 'CustomerCCNumber',
            c.customer_address AS 'CustomerAddress',
            c.customer_email AS 'CustomerEmail',
            c.customer_phone AS 'CustomerPhone',
            c.customer_tax_id AS 'CustomerTaxID',
            s.sales_channel AS 'SalesChannel',
            s.payment_method AS 'PaymentMethod',
            s.payment_status AS 'PaymentStatus',
            s.delivery_status AS 'DeliveryStatus',
            s.canceled AS 'IsCanceled',
            s.created_at AS 'OrderDate',
            s.payment_date AS 'PaymentDate',
            s.delivery_date AS 'DeliveryDate'
        FROM sales s
        JOIN customers c ON s.customer_id = c.customer_id
        JOIN inventory i ON s.item_id = i.inventory_id
        WHERE i.sellable = 1
        ORDER BY s.created_at DESC;
        `)

    if err != nil {
        return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't find sales", "data": nil})
    }
    defer salesRows.Close()

    salesOrdersMap := make(map[string]models.OptimizedGetSaleSensitive)

    for salesRows.Next() {
        var (
            salesModel     models.OptimizedGetSaleSensitive
            salesItemModel models.OptimizedGetSalesProduct
            customerModel  models.GetCustomerSensitive
        )

        err := salesRows.Scan(
            &salesModel.SalesOrderNumber,
            &salesItemModel.ProductID,
            &salesItemModel.ProductName,
            &salesItemModel.Quantity,
            &salesItemModel.SentQuantity,
            &salesItemModel.MfgPricePerUnit,
            &salesItemModel.SalePriceperUnit,
            &salesModel.TaxPercent,
            &customerModel.CustomerID,
            &customerModel.CustomerName,
            &customerModel.CustomerBankName,
            &customerModel.CustomerBankAccountNumber,
            &customerModel.CustomerCCNumber,
            &customerModel.CustomerAddress,
            &customerModel.CustomerEmail,
            &customerModel.CustomerPhone,
            &customerModel.CustomerTaxID,
            &salesModel.SalesChannel,
            &salesModel.PaymentMethod,
            &salesModel.PaymentStatus,
            &salesModel.DeliveryStatus,
            &salesModel.IsCanceled,
            &salesModel.OrderDate,
            &salesModel.PaymentDate,
            &salesModel.DeliveryDate,
        )

        if err != nil {
            return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't scan sales data", "data": nil})
        }

        // Check if the sales order already exists in the map
        if existingSale, exists := salesOrdersMap[salesModel.SalesOrderNumber]; exists {
            // If the sales order exists, just append the product to the list
            existingSale.Products = append(existingSale.Products, salesItemModel)

            // Prioritize Delivered Partially
            if salesModel.DeliveryStatus == "Delivered Partially" && existingSale.DeliveryStatus != "Delivered Partially" {
                existingSale.DeliveryStatus = "Delivered Partially"
            }

            salesOrdersMap[salesModel.SalesOrderNumber] = existingSale
        } else {
            // If the sales order doesn't exist, create a new one and append the product
            salesModel.Products = append(salesModel.Products, salesItemModel)
            salesModel.Customer = customerModel
            salesOrdersMap[salesModel.SalesOrderNumber] = salesModel
        }
    }

    // Convert map to slice for JSON response
    var salesModels []models.OptimizedGetSaleSensitive
    for _, sale := range salesOrdersMap {
        salesModels = append(salesModels, sale)
    }

    return c.JSON(salesModels)
}


func GetSales(c *fiber.Ctx) error {
	salesRows, err := database.DB.Query(
		`
		SELECT
			s.sales_order_number AS 'SalesOrderNumber',
			i.item_name AS 'ProductName',
			s.quantity AS 'Quantity',
			s.sent_quantity AS 'SentQuantity',
			i.price AS 'MfgPricePerUnit',
			s.sale_price_per_unit AS 'SalePriceperUnit',
			s.tax_percent AS 'TaxPercent',
			c.customer_id AS 'CustomerID',
			c.customer_name AS 'CustomerName',
			s.sales_channel AS 'SalesChannel',
			s.payment_method AS 'PaymentMethod',
			s.payment_status AS 'PaymentStatus',
			s.delivery_status AS 'DeliveryStatus',
			s.canceled AS 'IsCanceled',
			s.created_at AS 'OrderDate',
			s.payment_date AS 'PaymentDate',
			s.delivery_date AS 'DeliveryDate'
		FROM sales s
		JOIN customers c ON s.customer_id = c.customer_id
		JOIN inventory i ON s.item_id = i.inventory_id
		ORDER BY s.created_at DESC;
		`)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't find sales", "data": err})
	}

	salesModels := []models.GetSaleSensitive{}

	for salesRows.Next() {
		salesModel := models.GetSaleSensitive{}
		customerModel := models.GetCustomerSensitive{}

		salesRows.Scan(
			&salesModel.SalesOrderNumber,
			&salesModel.ProductName,
			&salesModel.Quantity,
			&salesModel.SentQuantity,
			&salesModel.MfgPricePerUnit,
			&salesModel.SalePriceperUnit,
			&salesModel.TaxPercent,
			&customerModel.CustomerID,
			&customerModel.CustomerName,
			&salesModel.SalesChannel,
			&salesModel.PaymentMethod,
			&salesModel.PaymentStatus,
			&salesModel.DeliveryStatus,
			&salesModel.IsCanceled,
			&salesModel.OrderDate,
			&salesModel.PaymentDate,
			&salesModel.DeliveryDate,
		)

		salesModel.Customer = customerModel
		salesModels = append(salesModels, salesModel)
	}

	defer salesRows.Close()

	return c.JSON(salesModels)
}

func GetCustomers(c *fiber.Ctx) error {
	customerModels := []models.GetCustomer{}
	
	customersRows, err := database.DB.Query(`SELECT customer_id,customer_name FROM customers;`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't find customers", "data": err})
	}

	for customersRows.Next() {
		customerModel := models.GetCustomer{}
		customersRows.Scan(&customerModel.CustomerID, &customerModel.CustomerName)
		customerModels = append(customerModels, customerModel)
	}

	defer customersRows.Close()

	return c.JSON(customerModels)
}

func EditSales(c *fiber.Ctx) error {
	salesOrderID := c.Params("sales_order_id")
	formData := c.Request().PostArgs()
	parsedForm := make(map[string][]string)

	formData.VisitAll(func(key, value []byte) {
		switch {
		case strings.HasPrefix(string(key), "product-id-"):
			parsedForm["product-id"] = append(parsedForm["product-id"], string(value))
		case strings.HasPrefix(string(key), "sent-quantity-"):
			parsedForm["sent-quantity"] = append(parsedForm["sent-quantity"], string(value))
		case strings.HasPrefix(string(key), "total-tax-"):
			parsedForm["total_tax"] = append(parsedForm["total_tax"], string(value))
		case strings.HasPrefix(string(key), "quantity-"):
			parsedForm["quantity"] = append(parsedForm["quantity"], string(value))
		case strings.HasPrefix(string(key), "sale-price-"):
			parsedForm["sale_price"] = append(parsedForm["sale_price"], string(value))
		case strings.Compare(string(key), "payment_status") == 0:
			parsedForm["payment_status"] = append(parsedForm["payment_status"], string(value))
		case strings.Compare(string(key), "payment_method") == 0:
			parsedForm["payment_method"] = append(parsedForm["payment_method"], string(value))
		case strings.Compare(string(key), "customer_id") == 0:
			parsedForm["customer_id"] = append(parsedForm["customer_id"], string(value))
		case strings.Compare(string(key), "sales_channel") == 0:
			parsedForm["sales_channel"] = append(parsedForm["sales_channel"], string(value))
		}
	})

	// debug
	// return c.JSON(parsedForm)

	// transaction
	tx, err := database.DB.Begin()
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	_, err = tx.Exec(`DELETE FROM sales WHERE sales_order_number = ?`, salesOrderID)
	if err != nil {
		tx.Rollback()
		return c.Status(500).SendString("Error at DELETE phase: " + err.Error())
	}

	numProductIDs := len(parsedForm["product-id"])

	var delivery_status string

	quantity, err := strconv.Atoi(parsedForm["quantity"][0])
	if err != nil {
		log.Println("Invalid quantity value:", err)
		delivery_status = "Pending"
	} 

	sentQuantity, err := strconv.Atoi(parsedForm["sent-quantity"][0])
	if err != nil {
		log.Println("Invalid sent-quantity value:", err)
		delivery_status = "Pending"
	}

	if quantity == sentQuantity {
		delivery_status = "Delivered"
	} else if sentQuantity > 0 && sentQuantity < quantity {
		delivery_status = "Delivered Partially"
	} else {
		delivery_status = "Pending"
	}

	query := `INSERT INTO sales (sales_order_number, item_id, quantity, sent_quantity, sale_price_per_unit, tax_percent, customer_id, sales_channel, payment_method, payment_status, delivery_status) VALUES `
	valueArgs := []interface{}{}
	var valueStrings []string

	paymentStatus := "Pending"
	if len(parsedForm["payment_status"]) > 0 {
		paymentStatus = parsedForm["payment_status"][0]
	}

	paymentMethod := "Cash"
	if len(parsedForm["payment_method"]) > 0 {
		paymentMethod = parsedForm["payment_method"][0]
	}

	for i := 0; i < numProductIDs; i++ {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, salesOrderID, parsedForm["product-id"][i], parsedForm["quantity"][i], parsedForm["sent-quantity"][i], parsedForm["sale_price"][i], parsedForm["total_tax"][0], parsedForm["customer_id"][0], parsedForm["sales_channel"][0], paymentMethod, paymentStatus, delivery_status)
	}

	query += strings.Join(valueStrings, ", ")

	// Execute the insert statement
	_, err = tx.Exec(query, valueArgs...)
	if err != nil {
		tx.Rollback()
		return c.Status(500).SendString("Error at INSERT phase: " + err.Error())
	}

	if err := tx.Commit(); err != nil {
		return c.Status(500).SendString("Error at COMMIT phase: " + err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "OK",
	})
}
