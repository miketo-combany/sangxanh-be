package database

import "database/sql"

type PublicUsersSelect struct {
  Address      []interface{}  `json:"address"`
  Avatar       sql.NullString `json:"avatar"`
  BasicAddress interface{}    `json:"basic_address"`
  CreatedAt    string         `json:"created_at"`
  DeletedAt    sql.NullString `json:"deleted_at"`
  Email        sql.NullString `json:"email"`
  Id           string         `json:"id"`
  Metadata     []interface{}  `json:"metadata"`
  Password     string         `json:"password"`
  Phone        sql.NullString `json:"phone"`
  Role         string         `json:"role"`
  UpdatedAt    sql.NullString `json:"updated_at"`
  Username     sql.NullString `json:"username"`
}

type PublicUsersInsert struct {
  Address      []interface{}  `json:"address"`
  Avatar       sql.NullString `json:"avatar"`
  BasicAddress interface{}    `json:"basic_address"`
  CreatedAt    sql.NullString `json:"created_at"`
  DeletedAt    sql.NullString `json:"deleted_at"`
  Email        sql.NullString `json:"email"`
  Id           sql.NullString `json:"id"`
  Metadata     []interface{}  `json:"metadata"`
  Password     string         `json:"password"`
  Phone        sql.NullString `json:"phone"`
  Role         sql.NullString `json:"role"`
  UpdatedAt    sql.NullString `json:"updated_at"`
  Username     sql.NullString `json:"username"`
}

type PublicUsersUpdate struct {
  Address      []interface{}  `json:"address"`
  Avatar       sql.NullString `json:"avatar"`
  BasicAddress interface{}    `json:"basic_address"`
  CreatedAt    sql.NullString `json:"created_at"`
  DeletedAt    sql.NullString `json:"deleted_at"`
  Email        sql.NullString `json:"email"`
  Id           sql.NullString `json:"id"`
  Metadata     []interface{}  `json:"metadata"`
  Password     sql.NullString `json:"password"`
  Phone        sql.NullString `json:"phone"`
  Role         sql.NullString `json:"role"`
  UpdatedAt    sql.NullString `json:"updated_at"`
  Username     sql.NullString `json:"username"`
}

type PublicCategoriesSelect struct {
  CreatedAt   string         `json:"created_at"`
  CreatedBy   sql.NullString `json:"created_by"`
  DeletedAt   sql.NullString `json:"deleted_at"`
  Description sql.NullString `json:"description"`
  Id          string         `json:"id"`
  Level       sql.NullInt32  `json:"level"`
  Metadata    interface{}    `json:"metadata"`
  Name        string         `json:"name"`
  ParentId    sql.NullString `json:"parent_id"`
  Status      sql.NullBool   `json:"status"`
  Thumbnail   sql.NullString `json:"thumbnail"`
  UpdatedAt   sql.NullString `json:"updated_at"`
}

type PublicCategoriesInsert struct {
  CreatedAt   sql.NullString `json:"created_at"`
  CreatedBy   sql.NullString `json:"created_by"`
  DeletedAt   sql.NullString `json:"deleted_at"`
  Description sql.NullString `json:"description"`
  Id          sql.NullString `json:"id"`
  Level       sql.NullInt32  `json:"level"`
  Metadata    interface{}    `json:"metadata"`
  Name        string         `json:"name"`
  ParentId    sql.NullString `json:"parent_id"`
  Status      sql.NullBool   `json:"status"`
  Thumbnail   sql.NullString `json:"thumbnail"`
  UpdatedAt   sql.NullString `json:"updated_at"`
}

type PublicCategoriesUpdate struct {
  CreatedAt   sql.NullString `json:"created_at"`
  CreatedBy   sql.NullString `json:"created_by"`
  DeletedAt   sql.NullString `json:"deleted_at"`
  Description sql.NullString `json:"description"`
  Id          sql.NullString `json:"id"`
  Level       sql.NullInt32  `json:"level"`
  Metadata    interface{}    `json:"metadata"`
  Name        sql.NullString `json:"name"`
  ParentId    sql.NullString `json:"parent_id"`
  Status      sql.NullBool   `json:"status"`
  Thumbnail   sql.NullString `json:"thumbnail"`
  UpdatedAt   sql.NullString `json:"updated_at"`
}

type PublicProductsSelect struct {
  CategoryId   sql.NullString  `json:"category_id"`
  Content      sql.NullString  `json:"content"`
  CreatedAt    string          `json:"created_at"`
  DeletedAt    sql.NullString  `json:"deleted_at"`
  Discount     sql.NullInt64   `json:"discount"`
  DiscountType sql.NullString  `json:"discount_type"`
  Id           string          `json:"id"`
  ImageDetail  interface{}     `json:"image_detail"`
  Metadata     []interface{}   `json:"metadata"`
  Name         string          `json:"name"`
  Price        sql.NullFloat64 `json:"price"`
  Thumbnail    sql.NullString  `json:"thumbnail"`
  UpdatedAt    sql.NullString  `json:"updated_at"`
}

type PublicProductsInsert struct {
  CategoryId   sql.NullString  `json:"category_id"`
  Content      sql.NullString  `json:"content"`
  CreatedAt    sql.NullString  `json:"created_at"`
  DeletedAt    sql.NullString  `json:"deleted_at"`
  Discount     sql.NullInt64   `json:"discount"`
  DiscountType sql.NullString  `json:"discount_type"`
  Id           sql.NullString  `json:"id"`
  ImageDetail  interface{}     `json:"image_detail"`
  Metadata     []interface{}   `json:"metadata"`
  Name         string          `json:"name"`
  Price        sql.NullFloat64 `json:"price"`
  Thumbnail    sql.NullString  `json:"thumbnail"`
  UpdatedAt    sql.NullString  `json:"updated_at"`
}

type PublicProductsUpdate struct {
  CategoryId   sql.NullString  `json:"category_id"`
  Content      sql.NullString  `json:"content"`
  CreatedAt    sql.NullString  `json:"created_at"`
  DeletedAt    sql.NullString  `json:"deleted_at"`
  Discount     sql.NullInt64   `json:"discount"`
  DiscountType sql.NullString  `json:"discount_type"`
  Id           sql.NullString  `json:"id"`
  ImageDetail  interface{}     `json:"image_detail"`
  Metadata     []interface{}   `json:"metadata"`
  Name         sql.NullString  `json:"name"`
  Price        sql.NullFloat64 `json:"price"`
  Thumbnail    sql.NullString  `json:"thumbnail"`
  UpdatedAt    sql.NullString  `json:"updated_at"`
}

type PublicProductVariantsSelect struct {
  CreatedAt string         `json:"created_at"`
  DeletedAt sql.NullString `json:"deleted_at"`
  Detail    interface{}    `json:"detail"`
  Id        string         `json:"id"`
  Metadata  []interface{}  `json:"metadata"`
  Name      string         `json:"name"`
  ProductId sql.NullString `json:"product_id"`
  UpdatedAt sql.NullString `json:"updated_at"`
}

type PublicProductVariantsInsert struct {
  CreatedAt sql.NullString `json:"created_at"`
  DeletedAt sql.NullString `json:"deleted_at"`
  Detail    interface{}    `json:"detail"`
  Id        sql.NullString `json:"id"`
  Metadata  []interface{}  `json:"metadata"`
  Name      string         `json:"name"`
  ProductId sql.NullString `json:"product_id"`
  UpdatedAt sql.NullString `json:"updated_at"`
}

type PublicProductVariantsUpdate struct {
  CreatedAt sql.NullString `json:"created_at"`
  DeletedAt sql.NullString `json:"deleted_at"`
  Detail    interface{}    `json:"detail"`
  Id        sql.NullString `json:"id"`
  Metadata  []interface{}  `json:"metadata"`
  Name      sql.NullString `json:"name"`
  ProductId sql.NullString `json:"product_id"`
  UpdatedAt sql.NullString `json:"updated_at"`
}

type PublicProductOptionsSelect struct {
  CreatedAt  sql.NullString  `json:"created_at"`
  DeletedAt  sql.NullString  `json:"deleted_at"`
  Detail     interface{}     `json:"detail"`
  ExtraPrice sql.NullFloat64 `json:"extra_price"`
  Id         string          `json:"id"`
  Image      sql.NullString  `json:"image"`
  Metadata   []interface{}   `json:"metadata"`
  Name       sql.NullString  `json:"name"`
  ProductId  sql.NullString  `json:"product_id"`
  Quantity   sql.NullInt32   `json:"quantity"`
  UpdatedAt  sql.NullString  `json:"updated_at"`
}

type PublicProductOptionsInsert struct {
  CreatedAt  sql.NullString  `json:"created_at"`
  DeletedAt  sql.NullString  `json:"deleted_at"`
  Detail     interface{}     `json:"detail"`
  ExtraPrice sql.NullFloat64 `json:"extra_price"`
  Id         sql.NullString  `json:"id"`
  Image      sql.NullString  `json:"image"`
  Metadata   []interface{}   `json:"metadata"`
  Name       sql.NullString  `json:"name"`
  ProductId  sql.NullString  `json:"product_id"`
  Quantity   sql.NullInt32   `json:"quantity"`
  UpdatedAt  sql.NullString  `json:"updated_at"`
}

type PublicProductOptionsUpdate struct {
  CreatedAt  sql.NullString  `json:"created_at"`
  DeletedAt  sql.NullString  `json:"deleted_at"`
  Detail     interface{}     `json:"detail"`
  ExtraPrice sql.NullFloat64 `json:"extra_price"`
  Id         sql.NullString  `json:"id"`
  Image      sql.NullString  `json:"image"`
  Metadata   []interface{}   `json:"metadata"`
  Name       sql.NullString  `json:"name"`
  ProductId  sql.NullString  `json:"product_id"`
  Quantity   sql.NullInt32   `json:"quantity"`
  UpdatedAt  sql.NullString  `json:"updated_at"`
}

type PublicOrdersSelect struct {
  Address   sql.NullString `json:"address"`
  CreatedAt sql.NullString `json:"created_at"`
  DeletedAt sql.NullString `json:"deleted_at"`
  Id        string         `json:"id"`
  Metadata  []interface{}  `json:"metadata"`
  Status    sql.NullString `json:"status"`
  UpdatedAt sql.NullString `json:"updated_at"`
  UserId    sql.NullString `json:"user_id"`
}

type PublicOrdersInsert struct {
  Address   sql.NullString `json:"address"`
  CreatedAt sql.NullString `json:"created_at"`
  DeletedAt sql.NullString `json:"deleted_at"`
  Id        sql.NullString `json:"id"`
  Metadata  []interface{}  `json:"metadata"`
  Status    sql.NullString `json:"status"`
  UpdatedAt sql.NullString `json:"updated_at"`
  UserId    sql.NullString `json:"user_id"`
}

type PublicOrdersUpdate struct {
  Address   sql.NullString `json:"address"`
  CreatedAt sql.NullString `json:"created_at"`
  DeletedAt sql.NullString `json:"deleted_at"`
  Id        sql.NullString `json:"id"`
  Metadata  []interface{}  `json:"metadata"`
  Status    sql.NullString `json:"status"`
  UpdatedAt sql.NullString `json:"updated_at"`
  UserId    sql.NullString `json:"user_id"`
}

type PublicOrderDetailsSelect struct {
  CreatedAt     sql.NullString `json:"created_at"`
  DeletedAt     sql.NullString `json:"deleted_at"`
  Discount      sql.NullInt32  `json:"discount"`
  DiscountType  sql.NullString `json:"discount_type"`
  Id            string         `json:"id"`
  Metadata      []interface{}  `json:"metadata"`
  OrderId       sql.NullString `json:"order_id"`
  ProductDetail []interface{}  `json:"product_detail"`
  Quantity      sql.NullInt32  `json:"quantity"`
  UpdatedAt     sql.NullString `json:"updated_at"`
}

type PublicOrderDetailsInsert struct {
  CreatedAt     sql.NullString `json:"created_at"`
  DeletedAt     sql.NullString `json:"deleted_at"`
  Discount      sql.NullInt32  `json:"discount"`
  DiscountType  sql.NullString `json:"discount_type"`
  Id            sql.NullString `json:"id"`
  Metadata      []interface{}  `json:"metadata"`
  OrderId       sql.NullString `json:"order_id"`
  ProductDetail []interface{}  `json:"product_detail"`
  Quantity      sql.NullInt32  `json:"quantity"`
  UpdatedAt     sql.NullString `json:"updated_at"`
}

type PublicOrderDetailsUpdate struct {
  CreatedAt     sql.NullString `json:"created_at"`
  DeletedAt     sql.NullString `json:"deleted_at"`
  Discount      sql.NullInt32  `json:"discount"`
  DiscountType  sql.NullString `json:"discount_type"`
  Id            sql.NullString `json:"id"`
  Metadata      []interface{}  `json:"metadata"`
  OrderId       sql.NullString `json:"order_id"`
  ProductDetail []interface{}  `json:"product_detail"`
  Quantity      sql.NullInt32  `json:"quantity"`
  UpdatedAt     sql.NullString `json:"updated_at"`
}

type PublicCartsSelect struct {
  CreatedAt     sql.NullString `json:"created_at"`
  DeletedAt     sql.NullString `json:"deleted_at"`
  Id            string         `json:"id"`
  Metadata      []interface{}  `json:"metadata"`
  ProductDetail []interface{}  `json:"product_detail"`
  Quantity      sql.NullInt32  `json:"quantity"`
  UpdatedAt     sql.NullString `json:"updated_at"`
  UserId        sql.NullString `json:"user_id"`
}

type PublicCartsInsert struct {
  CreatedAt     sql.NullString `json:"created_at"`
  DeletedAt     sql.NullString `json:"deleted_at"`
  Id            sql.NullString `json:"id"`
  Metadata      []interface{}  `json:"metadata"`
  ProductDetail []interface{}  `json:"product_detail"`
  Quantity      sql.NullInt32  `json:"quantity"`
  UpdatedAt     sql.NullString `json:"updated_at"`
  UserId        sql.NullString `json:"user_id"`
}

type PublicCartsUpdate struct {
  CreatedAt     sql.NullString `json:"created_at"`
  DeletedAt     sql.NullString `json:"deleted_at"`
  Id            sql.NullString `json:"id"`
  Metadata      []interface{}  `json:"metadata"`
  ProductDetail []interface{}  `json:"product_detail"`
  Quantity      sql.NullInt32  `json:"quantity"`
  UpdatedAt     sql.NullString `json:"updated_at"`
  UserId        sql.NullString `json:"user_id"`
}

type PublicPostsSelect struct {
  Assignee  sql.NullString `json:"assignee"`
  Content   sql.NullString `json:"content"`
  CreatedAt sql.NullString `json:"created_at"`
  CreatedBy sql.NullString `json:"created_by"`
  DeletedAt sql.NullString `json:"deleted_at"`
  Id        string         `json:"id"`
  Metadata  []interface{}  `json:"metadata"`
  Status    bool           `json:"status"`
  Thumbnail sql.NullString `json:"thumbnail"`
  Type      sql.NullString `json:"type"`
  UpdatedAt sql.NullString `json:"updated_at"`
}

type PublicPostsInsert struct {
  Assignee  sql.NullString `json:"assignee"`
  Content   sql.NullString `json:"content"`
  CreatedAt sql.NullString `json:"created_at"`
  CreatedBy sql.NullString `json:"created_by"`
  DeletedAt sql.NullString `json:"deleted_at"`
  Id        sql.NullString `json:"id"`
  Metadata  []interface{}  `json:"metadata"`
  Status    sql.NullBool   `json:"status"`
  Thumbnail sql.NullString `json:"thumbnail"`
  Type      sql.NullString `json:"type"`
  UpdatedAt sql.NullString `json:"updated_at"`
}

type PublicPostsUpdate struct {
  Assignee  sql.NullString `json:"assignee"`
  Content   sql.NullString `json:"content"`
  CreatedAt sql.NullString `json:"created_at"`
  CreatedBy sql.NullString `json:"created_by"`
  DeletedAt sql.NullString `json:"deleted_at"`
  Id        sql.NullString `json:"id"`
  Metadata  []interface{}  `json:"metadata"`
  Status    sql.NullBool   `json:"status"`
  Thumbnail sql.NullString `json:"thumbnail"`
  Type      sql.NullString `json:"type"`
  UpdatedAt sql.NullString `json:"updated_at"`
}

type PublicAuditTrailsSelect struct {
  AuditContent []interface{}  `json:"audit_content"`
  AuditId      sql.NullString `json:"audit_id"`
  AuditType    sql.NullString `json:"audit_type"`
  CreatedAt    sql.NullString `json:"created_at"`
  CreatedBy    sql.NullString `json:"created_by"`
  DeletedAt    sql.NullString `json:"deleted_at"`
  Id           string         `json:"id"`
  Metadata     []interface{}  `json:"metadata"`
  UpdatedAt    sql.NullString `json:"updated_at"`
}

type PublicAuditTrailsInsert struct {
  AuditContent []interface{}  `json:"audit_content"`
  AuditId      sql.NullString `json:"audit_id"`
  AuditType    sql.NullString `json:"audit_type"`
  CreatedAt    sql.NullString `json:"created_at"`
  CreatedBy    sql.NullString `json:"created_by"`
  DeletedAt    sql.NullString `json:"deleted_at"`
  Id           sql.NullString `json:"id"`
  Metadata     []interface{}  `json:"metadata"`
  UpdatedAt    sql.NullString `json:"updated_at"`
}

type PublicAuditTrailsUpdate struct {
  AuditContent []interface{}  `json:"audit_content"`
  AuditId      sql.NullString `json:"audit_id"`
  AuditType    sql.NullString `json:"audit_type"`
  CreatedAt    sql.NullString `json:"created_at"`
  CreatedBy    sql.NullString `json:"created_by"`
  DeletedAt    sql.NullString `json:"deleted_at"`
  Id           sql.NullString `json:"id"`
  Metadata     []interface{}  `json:"metadata"`
  UpdatedAt    sql.NullString `json:"updated_at"`
}
