package orm

import (
	"time"

	"gorm.io/gorm"
)

type Configuracion struct {
	gorm.Model
	EntregaTime      time.Time
	CambiosTime      time.Time
	PrecioSubvencion float32
}

type Usuario struct {
	gorm.Model
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
	gorm.Model
	Nombre string `gorm:"size:250"`
}

type Alergeno struct {
	gorm.Model
	Nombre string `gorm:"size:250"`
}

type Promocion struct {
	gorm.Model
	PlatoID uint // FK - promocion pertenece a un plato
	Inicio  time.Time
	Fin     time.Time
	Precio  float32 `gorm:"scale:2"`
}

type Plato struct {
	gorm.Model
	Nombre       string        `gorm:"unique;size:250"`
	Descripcion  string        `gorm:"size:2000"`
	Ingredientes []Ingrediente `gorm:"many2many:plato_ingredientes;"`
	Alergenos    []Alergeno    `gorm:"many2many:plato_alergenos;"`
	Precio       float32       `gorm:"scale:2"`
	Promociones  []Promocion   // has many
}

type PedidoLinea struct {
	gorm.Model
	PedidoID     uint    // FK - linea pertenece a un pedido
	Nombre       string  `gorm:"size:250"` // no usamos platos, nombre, precio, ingredientes podria cambiar
	PrecioUnidad float32 `gorm:"scale:2"`
	Cantidad     uint
}

type Pedido struct {
	gorm.Model
	PedidoLineas []PedidoLinea
	UsuarioID    uint    // FK - pedido pertenece a un usuario
	PrecioTotal  float32 `gorm:"scale:2"`
	PrecioPagar  float32 `gorm:"scale:2"` // FK - precio a pagar tras subvenciones
	Entrega      time.Time
}

type CarritoLinea struct {
	gorm.Model
	CarritoID uint // FK - linea pertenece a un carrito
	PlatoID   uint // FK - linea tiene un plato
}
type Carrito struct {
	gorm.Model
	CarritoLineas []CarritoLinea
	UsuarioID     int // FK - carrito pertenece a un usuario
}
