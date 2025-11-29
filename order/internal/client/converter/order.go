package converter

import (
	"google.golang.org/protobuf/types/known/wrapperspb"

	modelService "github.com/Bladforceone/rocket/order/internal/model"
	inventoryv1 "github.com/Bladforceone/rocket/shared/pkg/proto/inventory/v1"
)

func OrderToListRequest(order *modelService.Order) *inventoryv1.ListPartsRequest {
	return &inventoryv1.ListPartsRequest{
		Filter: &inventoryv1.PartFilter{
			Uuids:                 partUUIDsToProto(order),
			Names:                 nil,
			Categories:            nil,
			ManufacturerCountries: nil,
			Tags:                  nil,
		},
	}
}

func PartsPriceToService(response *inventoryv1.ListPartsResponse) []float64 {
	var partsPrice []float64
	for _, part := range response.Parts {
		partsPrice = append(partsPrice, part.GetPrice())
	}

	return partsPrice
}

func partUUIDsToProto(order *modelService.Order) []*wrapperspb.StringValue {
	var data []*wrapperspb.StringValue
	for _, v := range order.PartUUIDs {
		data = append(data, &wrapperspb.StringValue{
			Value: v,
		})
	}
	return data
}
