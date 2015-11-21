package srapi

import "testing"

func TestGetRegion(t *testing.T) {
	region, err := RegionById("mol4z19n")
	if err != nil {
		t.Fatal(err)
	}

	if len(region.Id) == 0 {
		t.Fatal("Region does not have an Id.")
	}

	if len(region.Name) == 0 {
		t.Fatal("Region does not have an name.")
	}
}

func TestGetNonexistingRegion(t *testing.T) {
	_, err := RegionById("i_do_not_exist")
	if err == nil {
		t.Fatal("Fetching a nonexisting region should yield an error.")
	}
}

func TestGetRegions(t *testing.T) {
	regions, err := Regions(nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(regions.Data) == 0 {
		t.Fatal("Could not find any regions.")
	}

	first := regions.Data[0]

	if len(first.Id) == 0 {
		t.Fatal("First region does not have an Id.")
	}

	if len(first.Name) == 0 {
		t.Fatal("First region does not have an name.")
	}
}

func TestGetRegionPages(t *testing.T) {
	regions, err := Regions(nil, &Cursor{0, 2})
	if err != nil {
		t.Fatal(err)
	}

	if regions.Pagination.Offset != 0 {
		t.Fatal("Requesting the first two regions did not return offset=0.")
	}

	prev, err := regions.PrevPage()

	if prev != nil {
		t.Fatal("Requesting the previous page of page 0 should yield nil.")
	}

	next, err := regions.NextPage()

	if next.Pagination.Offset != 2 {
		t.Fatal("The second page should have offset 2.")
	}
}
