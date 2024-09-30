package handlers

import (
	"rizaldyaristyo-fiber-boiler/database"
	"rizaldyaristyo-fiber-boiler/models"

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

		salesModel.Customer = &[]models.GetCustomerSensitive{customerModel}
		salesModels = append(salesModels, salesModel)
	}

	defer salesRows.Close()

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

		salesModel.Customer = &[]models.GetCustomerSensitive{customerModel}
		salesModels = append(salesModels, salesModel)
	}

	defer salesRows.Close()

	return c.JSON(salesModels)
}