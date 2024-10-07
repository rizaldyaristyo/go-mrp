package models

import (
	"database/sql"
	"time"
)

type User struct {
    UserID        string       `json:"user_id"`          // Custom employee ID
    Username      string       `json:"username"`
    Password      string       `json:"password"`
    Email         string       `json:"email"`
    Archived      bool         `json:"archived"`
    PermissionVal string       `json:"permission_val"`
    // PermissionVal
    //
    // 0 = Not Permitted
    // 1 = View
    // 2 = Add
    // 3 = Approve
    //
    // 0 0 0 0
    // │ │ │ └── Inventory
    // │ │ └──── Manufacturing
    // │ └────── Purchasing
    // └──────── Sales
    //
    // PermissionVal Example:
    // 1101 = View Sales, View Purchasing, No Manufacturing, View Inventory
    // 1211 = View Sales, Add Purchasing, View Manufacturing, View Inventory
    // 0312 = No Sales, Add Purchasing, View Manufacturing, View Inventory
}

type UserLogin struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type GetInventory struct {
    ManufacturingDemandQuantity sql.NullInt64   `json:"manufacturing_demand_quantity"`
    SalesDemandQuantity         sql.NullInt64   `json:"sales_demand_quantity"`
    TotalDemandQuantity         sql.NullInt64   `json:"total_demand_quantity"`
    InventoryID                 int             `json:"inventory_id"`
    ItemName                    string          `json:"item_name"`
    VendorID                    int             `json:"vendor_id"`
    VendorName                  string          `json:"vendor_name"`
    VendorAddress               string          `json:"vendor_address"`
    TaxID                       string          `json:"tax_id"`
    ItemCode                    string          `json:"item_code"`
    ItemCode2                   sql.NullString  `json:"item_code_2"`
    ItemType                    sql.NullString  `json:"item_type"`
    Sellable                    bool            `json:"sellable"`
    Purchaseable                bool            `json:"purchaseable"`
    Manufacturable              bool            `json:"manufacturable"`
    Price                       float64         `json:"price"`
    Price2                      sql.NullFloat64 `json:"price_2"`
    Currency                    sql.NullString  `json:"currency"`
    Quantity                    int             `json:"quantity"`
    MinimumStock                int             `json:"minimum_stock_warning"`
    LastUpdated                 time.Time       `json:"last_updated"`
    Archived                    bool            `json:"archived"`
    RecommendedMfgPrice         sql.NullFloat64 `json:"recommended_mfg_price"`
}

type GetManufacturingOrder struct {
    OrderID                 int                     `json:"order_id"`
    ManufactureOrderNumber  string                  `json:"manufacture_order_number"`
    ProductID               int                     `json:"product_id"`
    ProductName             string                  `json:"product_name"`
    QuantityToManufacture   int                     `json:"quantity_to_manufacture"`
    Status                  string                  `json:"status"`
    CreatedAt               sql.NullTime            `json:"created_at"`
    Archived                bool                    `json:"archived"`
    Recipes                 []ManufacturingRecipe   `json:"recipes"`
}

type ManufacturingRecipe struct {
    OrderID                         int       `json:"order_id"`
    ManufactureOrderNumber          string    `json:"manufacture_order_number"`
    ProductID                       int       `json:"product_id"`
    ProductName                     string    `json:"product_name"`
    QuantityToManufacture           int       `json:"quantity_to_manufacture"`
    Status                          string    `json:"status"`
    CreatedAt                       time.Time `json:"created_at"`
    Archived                        bool      `json:"archived"`
    MaterialQuantityToProduceOne    int       `json:"material_quantity_to_produce_one"`
    MaterialID                      int       `json:"material_id"`
    MaterialName                    string    `json:"material_name"`
    TotalMaterialQuantityNeeded     int       `json:"total_material_quantity_needed"`
    MaterialCurrentQuantity         int       `json:"material_current_quantity"`
}

type OptimizedGetManufacturingOrder struct {
    OrderID                 int                             `json:"order_id"`
    ManufactureOrderNumber  string                          `json:"manufacture_order_number"`
    ProductID               int                             `json:"product_id"`
    ProductName             string                          `json:"product_name"`
    QuantityToManufacture   int                             `json:"quantity_to_manufacture"`
    Status                  string                          `json:"status"`
    CreatedAt               sql.NullTime                    `json:"created_at"`
    Archived                bool                            `json:"archived"`
    Recipes                 *[]OptimizedManufacturingRecipe `json:"recipes"`
}

type OptimizedManufacturingRecipe struct {
    MaterialQuantityToProduceOne    sql.NullInt64       `json:"material_quantity_to_produce_one"`
    MaterialID                      sql.NullInt64       `json:"material_id"`
    MaterialName                    sql.NullString      `json:"material_name"`
    TotalMaterialQuantityNeeded     sql.NullInt64       `json:"total_material_quantity_needed"`
    MaterialCurrentQuantity         sql.NullInt64       `json:"material_current_quantity"`
}

type MaterialSufficiencies []MaterialSufficiency

type MaterialSufficiency struct {
    MaterialID      int       `json:"material_id"`
    MaterialName    string    `json:"material_name"`
    CurrentQuantity int       `json:"current_quantity"`
    NeededQuantity  int       `json:"needed_quantity"`
}

type GetVendors []GetVendor

type GetVendor struct {
    VendorID        int    `json:"vendor_id"`
    VendorName      string `json:"vendor_name"`
    VendorAddress   string `json:"vendor_address"`
    TaxID           string `json:"tax_id"`
}

type GetProducts struct{
    ProductName             string              `json:"product_name"`
    ProductInventoryCode    string              `json:"inventory_code"`
    ProductID               int                 `json:"product_id"`
    Recipes                 *[]GetProductRecipe `json:"recipes"`
}

type GetProductRecipe struct {
    MaterialID                      sql.NullInt64 `json:"material_id"`
    MaterialName                    sql.NullString `json:"material_name"`
    MaterialInventoryCode           sql.NullString `json:"material_inventory_code"`
    MaterialQuantityToProduceOne    sql.NullInt64 `json:"material_quantity_to_produce_one"`
    MaterialCurrentQuantity         sql.NullInt64 `json:"material_current_quantity"`
}

type GetMaterials struct {
    InventoryID        int     `json:"inventory_id"`
    ItemName           string  `json:"item_name"`
    ItemCode           string  `json:"item_code"`
}

type GetSaleSensitive struct {
    SalesOrderNumber   string                   `json:"sales_order_number"`
    ProductName        string                   `json:"product_name"`
    Quantity           int                      `json:"quantity"`
    SentQuantity       int                      `json:"sent_quantity"`
    MfgPricePerUnit    float64                  `json:"mfg_price_per_unit"`
    SalePriceperUnit   float64                  `json:"sale_price_per_unit"`
    TaxPercent         float64                  `json:"tax_percent"`
    Customer           GetCustomerSensitive     `json:"customer"`
    SalesChannel       string                   `json:"sales_channel"`
    PaymentMethod      string                   `json:"payment_method"`
    PaymentStatus      string                   `json:"payment_status"`
    DeliveryStatus     string                   `json:"delivery_status"`
    IsCanceled         bool                     `json:"is_canceled"`
    OrderDate          time.Time                `json:"order_date"`
    PaymentDate        sql.NullTime             `json:"payment_date"`
    DeliveryDate       sql.NullTime             `json:"delivery_date"`
}

//=======================================================================================

type GetCustomerSensitive struct {
    CustomerID                  int             `json:"customer_id"`
    CustomerName                string          `json:"customer_name"`
    CustomerBankName            sql.NullString  `json:"customer_bank_name"`
    CustomerBankAccountNumber   sql.NullString  `json:"customer_bank_account_number"`
    CustomerCCNumber            sql.NullString  `json:"customer_cc_number"`
    CustomerAddress             string          `json:"customer_address"`
    CustomerEmail               sql.NullString  `json:"customer_email"`
    CustomerPhone               sql.NullString  `json:"customer_phone"`
    CustomerTaxID               sql.NullString  `json:"customer_tax_id"`
}

type OptimizedGetSalesProduct struct {
    ProductID          string                   `json:"product_id"`
    ProductName        string                   `json:"product_name"`
    Quantity           int                      `json:"quantity"`
    SentQuantity       int                      `json:"sent_quantity"`
    MfgPricePerUnit    float64                  `json:"mfg_price_per_unit"`
    SalePriceperUnit   float64                  `json:"sale_price_per_unit"`
}

type OptimizedGetSaleSensitive struct {
    SalesOrderNumber   string                   `json:"sales_order_number"`
    Products           []OptimizedGetSalesProduct  `json:"products"`
    TaxPercent         float64                  `json:"tax_percent"`
    Customer           GetCustomerSensitive     `json:"customer"`
    SalesChannel       string                   `json:"sales_channel"`
    PaymentMethod      string                   `json:"payment_method"`
    PaymentStatus      string                   `json:"payment_status"`
    DeliveryStatus     string                   `json:"delivery_status"`
    IsCanceled         bool                     `json:"is_canceled"`
    OrderDate          time.Time                `json:"order_date"`
    PaymentDate        sql.NullTime             `json:"payment_date"`
    DeliveryDate       sql.NullTime             `json:"delivery_date"`
}

// ======================================================================================== //

type GetCustomer struct {
    CustomerID      int     `json:"customer_id"`
    CustomerName    string  `json:"customer_name"`
}

type GetSale struct {
    SalesOrderNumber   string           `json:"sales_order_number"`
    ProductName        string           `json:"product_name"`
    Quantity           int              `json:"quantity"`
    SentQuantity       int              `json:"sent_quantity"`
    MfgPricePerUnit    float64          `json:"mfg_price_per_unit"`
    SalePriceperUnit   float64          `json:"sale_price_per_unit"`
    TaxPercent         float64          `json:"tax_percent"`
    Customer           GetCustomer      `json:"customer"`
    SalesChannel       string           `json:"sales_channel"`
    PaymentMethod      string           `json:"payment_method"`
    PaymentStatus      string           `json:"payment_status"`
    DeliveryStatus     string           `json:"delivery_status"`
    IsCanceled         bool             `json:"is_canceled"`
    OrderDate          time.Time        `json:"order_date"`
    PaymentDate        sql.NullTime     `json:"payment_date"`
    DeliveryDate       sql.NullTime     `json:"delivery_date"`
}
