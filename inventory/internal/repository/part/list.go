package part

import (
	"context"

	"github.com/Bladforceone/rocket/inventory/internal/model"
	"github.com/Bladforceone/rocket/inventory/internal/repository/converter"
	repoModel "github.com/Bladforceone/rocket/inventory/internal/repository/model"
)

func (r *repository) List(ctx context.Context, filter model.PartFilter) ([]*model.Part, error) {
	fltr := converter.ToRepoPartFilter(&filter)

	r.mtx.RLock()
	defer r.mtx.RUnlock()

	var result []repoModel.Part

	for _, part := range r.data {
		if !uuidMatches(part.UUID, fltr.UUIDs) {
			continue
		}
		if !nameMatches(part.Name, fltr.Names) {
			continue
		}
		if !categoryMatches(part.Category, fltr.Categories) {
			continue
		}
		if !manufacturerMatches(part.Manufacturer, fltr.Manufacturers) {
			continue
		}
		if !tagsMatches(part.Tags, fltr.Tags) {
			continue
		}

		result = append(result, part)
	}

	return converter.ToServiceParts(&result), nil
}

func uuidMatches(uuid string, uuids []string) bool {
	if len(uuids) == 0 {
		return true
	}
	for _, u := range uuids {
		if uuid == u {
			return true
		}
	}
	return false
}

func nameMatches(name string, names []string) bool {
	if len(names) == 0 {
		return true
	}
	for _, n := range names {
		if name == n {
			return true
		}
	}
	return false
}

func categoryMatches(cat repoModel.Category, categories []repoModel.Category) bool {
	if len(categories) == 0 {
		return true
	}
	for _, c := range categories {
		if cat == c {
			return true
		}
	}
	return false
}

func manufacturerMatches(m repoModel.Manufacturer, manufacturers []repoModel.Manufacturer) bool {
	if len(manufacturers) == 0 {
		return true
	}
	for _, mf := range manufacturers {
		if m.Country == mf.Country {
			return true
		}
	}
	return false
}

func tagsMatches(partTags, filterTags []string) bool {
	if len(filterTags) == 0 {
		return true
	}
	for _, ft := range filterTags {
		for _, pt := range partTags {
			if ft == pt {
				return true
			}
		}
	}
	return false
}
