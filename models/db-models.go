package models

import (
	"database/sql"
	"time"
)

// 1. User Management and Permissions

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

// 2. Inventory

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
}


// ======================================================================================== //



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

type Sale struct {
    SaleID           int       `json:"sale_id"` 
    SalesOrderNumber string    `json:"sales_order_number"` 
    ItemID           int       `json:"item_id"`
    Quantity         int       `json:"quantity"` 
    SalePricePerUnit float64   `json:"sale_price_per_unit"` 
    SaleStatus       string    `json:"sale_status"` 
    SaleDate         time.Time `json:"sale_date"` 
    Archived         bool      `json:"archived"`
}

type Purchase struct {
    PurchaseID           int       `json:"purchase_id"` 
    PurchaseOrderNumber  string    `json:"purchase_order_number"` 
    ItemID               int       `json:"item_id"`
    Quantity             int       `json:"quantity"` 
    PurchasePricePerUnit float64   `json:"purchase_price_per_unit"` 
    PurchaseStatus       string    `json:"purchase_status"` 
    PurchaseDate         time.Time `json:"purchase_date"` 
    Archived             bool      `json:"archived"`
}

type Income struct {
    IncomeID     int       `json:"income_id"` 
    SaleID       int       `json:"sale_id"`
    Amount       float64   `json:"amount"` 
    Currency     string    `json:"currency"` 
    ReceivedDate time.Time `json:"received_date"` 
    Archived     bool      `json:"archived"`
}

type Outcome struct {
    OutcomeID    int       `json:"outcome_id"` 
    PurchaseID   int       `json:"purchase_id"`
    Amount       float64   `json:"amount"` 
    Currency     string    `json:"currency"` 
    SpentDate    time.Time `json:"spent_date"` 
    Archived     bool      `json:"archived"`
}