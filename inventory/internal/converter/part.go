package converter

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	serviceModel "github.com/Bladforceone/rocket/inventory/internal/model"
	inventoryv1 "github.com/Bladforceone/rocket/shared/pkg/proto/inventory/v1"
)

// ToServicePart конвертирует protobuf Part → service Part
func ToServicePart(p *inventoryv1.Part) *serviceModel.Part {
	if p == nil {
		return nil
	}

	return &serviceModel.Part{
		UUID:          p.Uuid,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		Category:      serviceModel.Category(p.Category),
		Dimensions:    ToServiceDimensions(p.Dimensions),
		Manufacturer:  ToServiceManufacturer(p.Manufacturer),
		Tags:          p.Tags,
		Metadata:      ToServiceMetadata(p.Metadata),
		CreatedAt:     ToTime(p.CreatedAt),
		UpdatedAt:     ToTime(p.UpdatedAt),
	}
}

// ToServiceParts repeated Part → []*serviceModel.Part
func ToServiceParts(parts []*inventoryv1.Part) []*serviceModel.Part {
	res := make([]*serviceModel.Part, 0, len(parts))
	for _, p := range parts {
		res = append(res, ToServicePart(p))
	}
	return res
}

// ToServiceDimensions конвертирует Dimensions
func ToServiceDimensions(d *inventoryv1.Dimensions) serviceModel.Dimensions {
	if d == nil {
		return serviceModel.Dimensions{}
	}
	return serviceModel.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

// ToServiceManufacturer конвертирует Manufacturer
func ToServiceManufacturer(m *inventoryv1.Manufacturer) serviceModel.Manufacturer {
	if m == nil {
		return serviceModel.Manufacturer{}
	}

	return serviceModel.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

// ToServiceMetadata map[string]*Value → map[string]interface{}
func ToServiceMetadata(meta map[string]*inventoryv1.Value) map[string]interface{} {
	if meta == nil {
		return nil
	}

	res := make(map[string]interface{}, len(meta))
	for k, v := range meta {
		switch val := v.Kind.(type) {
		case *inventoryv1.Value_StringValue:
			res[k] = val.StringValue
		case *inventoryv1.Value_IntValue:
			res[k] = val.IntValue
		case *inventoryv1.Value_DoubleValue:
			res[k] = val.DoubleValue
		case *inventoryv1.Value_BoolValue:
			res[k] = val.BoolValue
		}
	}
	return res
}

// ToServicePartFilter protobuf → service
func ToServicePartFilter(f *inventoryv1.PartFilter) *serviceModel.PartFilter {
	if f == nil {
		return nil
	}

	categories := make([]serviceModel.Category, len(f.Categories))
	for i, c := range f.Categories {
		categories[i] = serviceModel.Category(c)
	}

	uuids := unwrapStringValues(f.Uuids)
	names := unwrapStringValues(f.Names)
	manufacturers := unwrapStringValues(f.ManufacturerCountries)
	tags := unwrapStringValues(f.Tags)

	return &serviceModel.PartFilter{
		UUIDs:                 uuids,
		Names:                 names,
		Categories:            categories,
		ManufacturerCountries: manufacturers,
		Tags:                  tags,
	}
}

func unwrapStringValues(values []*wrapperspb.StringValue) []string {
	res := make([]string, 0, len(values))
	for _, v := range values {
		if v != nil {
			res = append(res, v.Value)
		}
	}
	return res
}

func ToTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}

// ToAPIPart конвертирует service Part → protobuf Part
func ToAPIPart(p *serviceModel.Part) *inventoryv1.Part {
	if p == nil {
		return nil
	}

	return &inventoryv1.Part{
		Uuid:          p.UUID,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		Category:      inventoryv1.Category(p.Category),
		Dimensions:    ToAPIDimensions(p.Dimensions),
		Manufacturer:  ToAPIManufacturer(p.Manufacturer),
		Tags:          append([]string(nil), p.Tags...),
		Metadata:      ToAPIMetadata(p.Metadata),
		CreatedAt:     timestamppb.New(p.CreatedAt),
		UpdatedAt:     timestamppb.New(p.UpdatedAt),
	}
}

// ToAPIParts []*serviceModel.Part → []*inventoryv1.Part
func ToAPIParts(parts []*serviceModel.Part) []*inventoryv1.Part {
	res := make([]*inventoryv1.Part, 0, len(parts))
	for _, p := range parts {
		res = append(res, ToAPIPart(p))
	}
	return res
}

// ToAPIDimensions service → protobuf
func ToAPIDimensions(d serviceModel.Dimensions) *inventoryv1.Dimensions {
	return &inventoryv1.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

// ToAPIManufacturer service → protobuf
func ToAPIManufacturer(m serviceModel.Manufacturer) *inventoryv1.Manufacturer {
	return &inventoryv1.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

// ToAPIMetadata map[string]interface{} → map[string]*inventoryv1.Value
func ToAPIMetadata(meta map[string]interface{}) map[string]*inventoryv1.Value {
	if meta == nil {
		return nil
	}

	res := make(map[string]*inventoryv1.Value, len(meta))
	for k, v := range meta {
		switch val := v.(type) {
		case string:
			res[k] = &inventoryv1.Value{Kind: &inventoryv1.Value_StringValue{StringValue: val}}
		case int64:
			res[k] = &inventoryv1.Value{Kind: &inventoryv1.Value_IntValue{IntValue: val}}
		case float64:
			res[k] = &inventoryv1.Value{Kind: &inventoryv1.Value_DoubleValue{DoubleValue: val}}
		case bool:
			res[k] = &inventoryv1.Value{Kind: &inventoryv1.Value_BoolValue{BoolValue: val}}
		}
	}
	return res
}

// ToAPIPartFilter service → protobuf
func ToAPIPartFilter(f *serviceModel.PartFilter) *inventoryv1.PartFilter {
	if f == nil {
		return nil
	}

	categories := make([]inventoryv1.Category, len(f.Categories))
	for i, c := range f.Categories {
		categories[i] = inventoryv1.Category(c)
	}

	return &inventoryv1.PartFilter{
		Uuids:                 wrapStringValues(f.UUIDs),
		Names:                 wrapStringValues(f.Names),
		Categories:            categories,
		ManufacturerCountries: wrapStringValues(f.ManufacturerCountries),
		Tags:                  wrapStringValues(f.Tags),
	}
}

func wrapStringValues(values []string) []*wrapperspb.StringValue {
	res := make([]*wrapperspb.StringValue, 0, len(values))
	for _, v := range values {
		res = append(res, &wrapperspb.StringValue{Value: v})
	}
	return res
}
