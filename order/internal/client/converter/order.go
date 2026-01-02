package converter

import (
	"google.golang.org/protobuf/types/known/wrapperspb"

	modelService "github.com/Bladforceone/rocket/order/internal/model"
	inventoryv1 "github.com/Bladforceone/rocket/shared/pkg/proto/inventory/v1"
)

func PartsToService(part *inventoryv1.ListPartsResponse) []modelService.Part {
	parts := make([]modelService.Part, len(part.Parts))
	for _, v := range part.Parts {
		parts = append(parts, PartToService(v))
	}
	return parts
}

func PartToService(part *inventoryv1.Part) modelService.Part {
	metadata := make(map[string]interface{})
	for k, v := range part.GetMetadata() {
		metadata[k] = v
	}

	return modelService.Part{
		UUID:          part.GetUuid(),
		Name:          part.GetName(),
		Description:   part.GetDescription(),
		Price:         part.GetPrice(),
		StockQuantity: part.GetStockQuantity(),
		Category:      modelService.Category(part.GetCategory()),
		Dimensions:    dimensionsToService(part.GetDimensions()),
		Manufacturer:  manufacturerToService(part.GetManufacturer()),
		Tags:          part.GetTags(),
		Metadata:      metadata,
		CreatedAt:     part.CreatedAt.AsTime(),
		UpdatedAt:     part.UpdatedAt.AsTime(),
	}
}

func dimensionsToService(dimens *inventoryv1.Dimensions) modelService.Dimensions {
	return modelService.Dimensions{
		Length: dimens.GetLength(),
		Width:  dimens.GetWidth(),
		Height: dimens.GetHeight(),
		Weight: dimens.GetWeight(),
	}
}

func manufacturerToService(manufact *inventoryv1.Manufacturer) modelService.Manufacturer {
	return modelService.Manufacturer{
		Name:    manufact.GetName(),
		Country: manufact.GetCountry(),
		Website: manufact.GetWebsite(),
	}
}

func FilterToListRequest(filter *modelService.PartFilter) *inventoryv1.ListPartsRequest {
	return &inventoryv1.ListPartsRequest{
		Filter: &inventoryv1.PartFilter{
			Uuids:                 stringsToProto(filter.UUIDs),
			Names:                 stringsToProto(filter.Names),
			Categories:            categoriesToProto(filter.Categories),
			ManufacturerCountries: stringsToProto(filter.ManufacturerCountries),
			Tags:                  stringsToProto(filter.Tags),
		},
	}
}

func stringsToProto(strs []string) []*wrapperspb.StringValue {
	var data []*wrapperspb.StringValue
	for _, v := range strs {
		data = append(data, &wrapperspb.StringValue{
			Value: v,
		})
	}
	return data
}

func categoriesToProto(categories []string) []inventoryv1.Category {
	var data []inventoryv1.Category
	for _, v := range categories {
		data = append(data, inventoryv1.Category(inventoryv1.Category_value[v]))
	}
	return data
}
