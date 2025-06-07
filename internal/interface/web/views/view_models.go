package views

// Aquí definimos los modelos o DTOs específicos para las vistas
// Estos modelos son los que se pasarán a las templates para renderizar

// UserViewModel representa los datos de usuario para las vistas
type UserViewModel struct {
	ID       uint
	Username string
	Email    string
	IsAdmin  bool
}

// ProductViewModel representa los datos de producto para las vistas
type ProductViewModel struct {
	ID             uint
	Name           string
	Slug           string
	Description    string
	ImageURL       string
	BestPrice      float64
	BestStore      string
	Category       CategoryViewModel
	Specifications []SpecificationViewModel
	Prices         []PriceViewModel
}

// CategoryViewModel representa los datos de categoría para las vistas
type CategoryViewModel struct {
	ID           uint
	Name         string
	Slug         string
	ProductCount int
}

// PriceViewModel representa los datos de precio para las vistas
type PriceViewModel struct {
	Store string
	Price float64
	URL   string
}

// SimilarProductViewModel representa un producto similar para mostrar en "Productos similares"
type SimilarProductViewModel struct {
	ID       uint
	Name     string
	ImageURL string
	Price    float64
	Store    string
	URL      string
}

// SpecificationViewModel representa una especificación técnica de un producto
type SpecificationViewModel struct {
	Name  string
	Value string
}

// HomePageViewModel representa el modelo para la vista de la página principal
type HomePageViewModel struct {
	User             *UserViewModel
	FeaturedProducts []ProductViewModel
	Categories       []CategoryViewModel
}

// ProductDetailViewModel representa el modelo para la vista de detalle de producto
type ProductDetailViewModel struct {
	User            *UserViewModel
	Product         ProductViewModel
	BestPrice       PriceViewModel
	OtherPrices     []PriceViewModel
	RelatedProducts []ProductViewModel
	LastUpdated     string
}

// CategoryPageViewModel representa el modelo para la vista de categoría
type CategoryPageViewModel struct {
	User      *UserViewModel
	Category  CategoryViewModel
	Products  []ProductViewModel
	PageCount int
	Page      int
}
