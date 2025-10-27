package converter

import (
	serviceModel "github.com/Bladforceone/rocket/inventory/internal/model"
	repoModel "github.com/Bladforceone/rocket/inventory/internal/repository/model"
)

// ToServicePart конвертирует модель из репозитория в сервисную
func ToServicePart(p *repoModel.Part) *serviceModel.Part {
	return &serviceModel.Part{
		UUID:          p.UUID,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		Category:      serviceModel.Category(p.Category),
		Dimensions:    serviceModel.Dimensions(p.Dimensions),
		Manufacturer:  serviceModel.Manufacturer(p.Manufacturer),
		Tags:          append([]string(nil), p.Tags...),
		Metadata:      p.Metadata,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

// ToServiceParts конвертирует слайс моделей из репозитория в слайс сервисных
func ToServiceParts(p *[]repoModel.Part) []*serviceModel.Part {
	var result []*serviceModel.Part
	for _, part := range *p {
		result = append(result, ToServicePart(&part))
	}

	return result
}

// ToRepoPartFilter конвертирует фильтр из сервисного слоя в репозиторий
func ToRepoPartFilter(f *serviceModel.PartFilter) *repoModel.PartFilter {
	categories := make([]repoModel.Category, 0, len(f.Categories))
	for _, category := range f.Categories {
		categories = append(categories, repoModel.Category(category))
	}

	return &repoModel.PartFilter{
		UUIDs:                 append([]string(nil), f.UUIDs...),
		Names:                 append([]string(nil), f.Names...),
		Categories:            categories,
		ManufacturerCountries: f.ManufacturerCountries,
		Tags:                  append([]string(nil), f.Tags...),
	}
}
