package util

import "testing"

func TestParsePagination(t *testing.T) {
    tests := []struct {
        name     string
        pageStr  string
        sizeStr  string
        wantPage int
        wantSize int
        wantErr  bool
    }{
        {"default values", "", "", 1, 10, false},
        {"custom page", "5", "", 5, 10, false},
        {"custom size", "", "20", 1, 20, false},
        {"both custom", "3", "25", 3, 25, false},
        {"invalid page", "abc", "", 0, 0, true},
        {"negative page", "-1", "", 0, 0, true},
        {"zero page", "0", "", 0, 0, true},
        {"size too large", "", "200", 0, 0, true},
        {"page with spaces", " 5 ", "", 5, 10, false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            page, size, err := ParsePagination(tt.pageStr, tt.sizeStr)

            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if !tt.wantErr {
                if page != tt.wantPage {
                    t.Errorf("page = %d, want %d", page, tt.wantPage)
                }
                if size != tt.wantSize {
                    t.Errorf("size = %d, want %d", size, tt.wantSize)
                }
            }
        })
    }
}
