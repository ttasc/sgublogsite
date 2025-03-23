package controller

import (
	"fmt"
	"net/url"
)

const postsLimitPerPage = 10
const maxVisiblePageItem = 5

type paginationItem struct {
    Number  int
    URL     string
    Active  bool
    Label   string
    BaseURL string
}

func generatePagination(baseURL string, currentPage, totalPages int) []paginationItem {
    var items []paginationItem

    u, _ := url.Parse(baseURL)
    query := u.Query()

    // Previous
    if currentPage > 1 {
        query.Set("page", fmt.Sprint(currentPage-1))
        u.RawQuery = query.Encode()
        items = append(items, paginationItem{
            Label:   "←",
            URL:     u.String(),
            Active:  false,
            BaseURL: baseURL,
        })
    }

    // Tính toán range trang
    start := max(currentPage-2, 1)
    end := min(start+maxVisiblePageItem-1, totalPages)

    for i := start; i <= end; i++ {
        query.Set("page", fmt.Sprint(i))
        u.RawQuery = query.Encode()
        items = append(items, paginationItem{
            Number:  i,
            URL:     u.String(),
            Active:  i == currentPage,
            BaseURL: baseURL,
        })
    }

    // Next
    if currentPage < totalPages {
        query.Set("page", fmt.Sprint(currentPage+1))
        u.RawQuery = query.Encode()
        items = append(items, paginationItem{
            Label:   "→",
            URL:     u.String(),
            Active:  false,
            BaseURL: baseURL,
        })
    }

    return items
}

