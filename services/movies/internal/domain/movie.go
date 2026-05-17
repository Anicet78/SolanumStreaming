package domain

type SearchRequestQuery struct {
	Title string `query:"title"`
}

type SearchResponse struct {
}

type TMDBSearchResponse struct {
	results any
}
