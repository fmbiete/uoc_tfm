package orm

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint64 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Configuracion struct {
	BaseModel
	EntregaTime      time.Time
	CambiosTime      time.Time
	PrecioSubvencion float64
}

type Usuario struct {
	BaseModel
	Email         string `gorm:"size:100;index:ix_login,priority:1"`
	Password      string `gorm:"size:250;index:ix_login,priority:2"`
	Nombre        string `gorm:"size:250"`
	Apellidos     string `gorm:"size:250"`
	Direccion1    string `gorm:"size:250"`
	Direccion2    string `gorm:"size:250"`
	Direccion3    string `gorm:"size:250"`
	Ciudad        string `gorm:"size:250"`
	CodigoPostal  string `gorm:"size:10"`
	Telefono      string `gorm:"size:20"`
	IsRestaurador bool
	Pedidos       []Pedido // has many
	CarritoCompra Carrito  // has one
}

type Ingrediente struct {
	BaseModel
	Nombre string `gorm:"size:250"`
}

type Alergeno struct {
	BaseModel
	Nombre string `gorm:"size:250"`
}

type Promocion struct {
	BaseModel
	PlatoID uint // FK - promocion pertenece a un plato
	Inicio  time.Time
	Fin     time.Time
	Precio  float64 `gorm:"scale:2"`
}

type Plato struct {
	BaseModel
	Nombre       string        `gorm:"unique;size:250"`
	Descripcion  string        `gorm:"size:2000"`
	Ingredientes []Ingrediente `gorm:"many2many:plato_ingredientes;"`
	Alergenos    []Alergeno    `gorm:"many2many:plato_alergenos;"`
	Precio       float64       `gorm:"scale:2"`
	Promociones  []Promocion   // has many
}

type PedidoLinea struct {
	BaseModel
	PedidoID     uint64  // FK - linea pertenece a un pedido
	Nombre       string  `gorm:"size:250"` // no usamos platos, nombre, precio, ingredientes podria cambiar
	PrecioUnidad float64 `gorm:"scale:2"`
	Cantidad     uint
}

type Pedido struct {
	BaseModel
	PedidoLineas []PedidoLinea
	UsuarioID    uint64  // FK - pedido pertenece a un usuario
	PrecioTotal  float64 `gorm:"scale:2"`
	PrecioPagar  float64 `gorm:"scale:2"` // FK - precio a pagar tras subvenciones
	Entrega      time.Time
}

type CarritoLinea struct {
	BaseModel
	CarritoID uint64 // FK - linea pertenece a un carrito
	PlatoID   uint64 // FK - linea tiene un plato
	Cantidad  uint
}

type Carrito struct {
	BaseModel
	CarritoLineas []CarritoLinea
	UsuarioID     uint64 // FK - carrito pertenece a un usuario
}
