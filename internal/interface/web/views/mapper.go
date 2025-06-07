package views

import (
	"time"

	"app/internal/domain/model"
)

// ToUserViewModel convierte un modelo de dominio User a un ViewModel para presentación
func ToUserViewModel(user *model.User) *UserViewModel {
	if user == nil {
		return nil
	}

	return &UserViewModel{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		IsAdmin:  user.IsAdmin,
	}
}

// ToProductViewModel convierte un modelo de dominio Product a un ViewModel
func ToProductViewModel(product *model.Product, bestPrice *model.Price) ProductViewModel {
	var categoryVM CategoryViewModel
	if product.CategoryID > 0 {
		categoryVM = ToCategoryViewModel(product.Category, 0)
	}

	// Crear especificaciones
	specs := make([]SpecificationViewModel, 0)
	for key, value := range product.Specifications {
		specs = append(specs, SpecificationViewModel{
			Name:  key,
			Value: value,
		})
	}

	bestPriceValue := 0.0
	bestStore := ""
	if bestPrice != nil {
		bestPriceValue = bestPrice.Price
		bestStore = bestPrice.Store
	}

	return ProductViewModel{
		ID:             product.ID,
		Name:           product.Name,
		Slug:           product.Slug,
		Description:    product.Description,
		ImageURL:       product.ImageURL,
		BestPrice:      bestPriceValue,
		BestStore:      bestStore,
		Category:       categoryVM,
		Specifications: specs,
	}
}

// ToCategoryViewModel convierte un modelo de dominio Category a ViewModel
func ToCategoryViewModel(category model.Category, productCount int) CategoryViewModel {
	return CategoryViewModel{
		ID:           category.ID,
		Name:         category.Name,
		Slug:         category.Slug,
		ProductCount: productCount,
	}
}

// ToPriceViewModel convierte un modelo de dominio Price a ViewModel
func ToPriceViewModel(price model.Price) PriceViewModel {
	return PriceViewModel{
		Store: price.Store,
		Price: price.Price,
		URL:   price.URL,
	}
}

// BuildHomePageViewModel construye un modelo completo para la vista de la página de inicio
func BuildHomePageViewModel(
	user *model.User,
	featuredProducts []*model.Product,
	categories []model.Category,
	bestPrices map[uint]*model.Price,
	productCounts map[uint]int,
) HomePageViewModel {

	// Construir los productos destacados
	productViewModels := make([]ProductViewModel, 0, len(featuredProducts))
	for _, product := range featuredProducts {
		bestPrice := bestPrices[product.ID]
		productViewModels = append(productViewModels, ToProductViewModel(product, bestPrice))
	}

	// Construir las categorías
	categoryViewModels := make([]CategoryViewModel, 0, len(categories))
	for _, category := range categories {
		count := productCounts[category.ID]
		categoryViewModels = append(categoryViewModels, ToCategoryViewModel(category, count))
	}

	return HomePageViewModel{
		User:             ToUserViewModel(user),
		FeaturedProducts: productViewModels,
		Categories:       categoryViewModels,
	}
}

// BuildProductDetailViewModel construye un modelo completo para la vista de detalle de producto
func BuildProductDetailViewModel(
	user *model.User,
	product *model.Product,
	prices []model.Price,
	relatedProducts []*model.Product,
	relatedBestPrices map[uint]*model.Price,
) ProductDetailViewModel {

	var bestPrice PriceViewModel
	otherPrices := make([]PriceViewModel, 0)

	// Obtener el mejor precio y los demás precios
	if len(prices) > 0 {
		bestPrice = ToPriceViewModel(prices[0])

		for i := 1; i < len(prices); i++ {
			otherPrices = append(otherPrices, ToPriceViewModel(prices[i]))
		}
	}

	// Productos relacionados
	relatedProductVMs := make([]ProductViewModel, 0, len(relatedProducts))
	for _, relProduct := range relatedProducts {
		bestPrice := relatedBestPrices[relProduct.ID]
		relatedProductVMs = append(relatedProductVMs, ToProductViewModel(relProduct, bestPrice))
	}

	// Fecha de última actualización
	lastUpdated := time.Now().Format("02/01/2006 15:04")

	return ProductDetailViewModel{
		User:            ToUserViewModel(user),
		Product:         ToProductViewModel(product, &prices[0]),
		BestPrice:       bestPrice,
		OtherPrices:     otherPrices,
		RelatedProducts: relatedProductVMs,
		LastUpdated:     lastUpdated,
	}
}
