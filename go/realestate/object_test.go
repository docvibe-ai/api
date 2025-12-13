package realestate

import (
	"testing"

	"github.com/domonda/go-types/country"
	"github.com/domonda/go-types/notnull"
)

func TestObject_Normalize(t *testing.T) {
	tests := []struct {
		name               string
		object             *Object
		wantErrors         int
		wantCountry        country.NullableCode
		wantStreetVarCount int
		wantStreetVars     []notnull.TrimmedString
	}{
		{
			name:       "nil returns no errors",
			object:     nil,
			wantErrors: 0,
		},
		{
			name:               "empty struct returns no errors",
			object:             &Object{},
			wantErrors:         0,
			wantCountry:        "",
			wantStreetVarCount: 0,
		},
		{
			name: "valid country code is normalized",
			object: &Object{
				ID:      "obj1",
				Country: "at",
			},
			wantErrors:  0,
			wantCountry: "AT",
		},
		{
			name: "invalid country code is cleared",
			object: &Object{
				ID:      "obj1",
				Country: "XX",
			},
			wantErrors:  1,
			wantCountry: "",
		},
		{
			name: "valid street variations are kept",
			object: &Object{
				ID:               "obj1",
				StreetVariations: []notnull.TrimmedString{"Main Street", "Main St"},
			},
			wantErrors:         0,
			wantStreetVarCount: 2,
			wantStreetVars:     []notnull.TrimmedString{"Main Street", "Main St"},
		},
		{
			name: "empty street variations are removed",
			object: &Object{
				ID:               "obj1",
				StreetVariations: []notnull.TrimmedString{"Main Street", "", "  ", "Main St"},
			},
			wantErrors:         0,
			wantStreetVarCount: 2,
			wantStreetVars:     []notnull.TrimmedString{"Main Street", "Main St"},
		},
		{
			name: "all fields normalized together",
			object: &Object{
				ID:               "obj1",
				Type:             "apartment",
				Street:           "Main Street 1",
				City:             "Vienna",
				Country:          "at",
				StreetVariations: []notnull.TrimmedString{"Main St 1", "", "Hauptstraße 1"},
			},
			wantErrors:         0,
			wantCountry:        "AT",
			wantStreetVarCount: 2,
			wantStreetVars:     []notnull.TrimmedString{"Main St 1", "Hauptstraße 1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.object.Normalize()
			if len(errs) != tt.wantErrors {
				t.Errorf("Normalize() got %d errors, want %d: %v", len(errs), tt.wantErrors, errs)
			}
			if tt.object == nil {
				return
			}
			if tt.object.Country != tt.wantCountry {
				t.Errorf("Country = %q, want %q", tt.object.Country, tt.wantCountry)
			}
			if len(tt.object.StreetVariations) != tt.wantStreetVarCount {
				t.Errorf("StreetVariations length = %d, want %d", len(tt.object.StreetVariations), tt.wantStreetVarCount)
			}
			if tt.wantStreetVars != nil {
				for i, v := range tt.wantStreetVars {
					if i < len(tt.object.StreetVariations) && tt.object.StreetVariations[i] != v {
						t.Errorf("StreetVariations[%d] = %q, want %q", i, tt.object.StreetVariations[i], v)
					}
				}
			}
		})
	}
}
