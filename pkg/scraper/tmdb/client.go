package tmdb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	tmdbAPIBase   = "https://api.themoviedb.org/3"
	tmdbImageBase = "https://image.tmdb.org/t/p/original"
)

type Client struct {
	apiKey string
	client *http.Client
}

// Common structures
type Credits struct {
	Cast []struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Character string `json:"character"`
	} `json:"cast"`
	Crew []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Job  string `json:"job"`
	} `json:"crew"`
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Movie structures
type MovieSearchResult struct {
	Page    int     `json:"page"`
	Results []Movie `json:"results"`
}

type Movie struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	ReleaseDate string  `json:"release_date"`
	PosterPath  string  `json:"poster_path"`
	Overview    string  `json:"overview"`
	VoteAverage float64 `json:"vote_average"`
	Runtime     int     `json:"runtime"`
	Genres      []Genre `json:"genres"`
}

type MovieDetails struct {
	Movie
	Credits Credits `json:"credits"`
}

// TV Show structures
type TVSearchResult struct {
	Page    int      `json:"page"`
	Results []TVShow `json:"results"`
}

type TVShow struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	FirstAirDate string  `json:"first_air_date"`
	PosterPath   string  `json:"poster_path"`
	Overview     string  `json:"overview"`
	VoteAverage  float64 `json:"vote_average"`
	Genres       []Genre `json:"genres"`
}

type TVShowDetails struct {
	TVShow
	Credits  Credits   `json:"credits"`
	Seasons  []Season  `json:"seasons"`
	Episodes []Episode `json:"episodes"`
}

type Season struct {
	ID           int    `json:"id"`
	SeasonNumber int    `json:"season_number"`
	Name         string `json:"name"`
	Overview     string `json:"overview"`
	PosterPath   string `json:"poster_path"`
}

type Episode struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Overview      string `json:"overview"`
	StillPath     string `json:"still_path"`
	AirDate       string `json:"air_date"`
	EpisodeNumber int    `json:"episode_number"`
	SeasonNumber  int    `json:"season_number"`
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Movie methods
func (c *Client) SearchMovie(query string) (*MovieSearchResult, error) {
	endpoint := fmt.Sprintf("%s/search/movie", tmdbAPIBase)
	params := url.Values{}
	params.Add("api_key", c.apiKey)
	params.Add("query", query)

	var result MovieSearchResult
	if err := c.get(endpoint, params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetMovieDetails(movieID int) (*MovieDetails, error) {
	endpoint := fmt.Sprintf("%s/movie/%d", tmdbAPIBase, movieID)
	params := url.Values{}
	params.Add("api_key", c.apiKey)
	params.Add("append_to_response", "credits")

	var details MovieDetails
	if err := c.get(endpoint, params, &details); err != nil {
		return nil, err
	}
	return &details, nil
}

// TV Show methods
func (c *Client) SearchTV(query string) (*TVSearchResult, error) {
	endpoint := fmt.Sprintf("%s/search/tv", tmdbAPIBase)
	params := url.Values{}
	params.Add("api_key", c.apiKey)
	params.Add("query", query)

	var result TVSearchResult
	if err := c.get(endpoint, params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetTVShowDetails(tvID int) (*TVShowDetails, error) {
	endpoint := fmt.Sprintf("%s/tv/%d", tmdbAPIBase, tvID)
	params := url.Values{}
	params.Add("api_key", c.apiKey)
	params.Add("append_to_response", "credits")

	var details TVShowDetails
	if err := c.get(endpoint, params, &details); err != nil {
		return nil, err
	}
	return &details, nil
}

func (c *Client) GetSeasonDetails(tvID int, seasonNumber int) (*Season, error) {
	endpoint := fmt.Sprintf("%s/tv/%d/season/%d", tmdbAPIBase, tvID, seasonNumber)
	params := url.Values{}
	params.Add("api_key", c.apiKey)

	var season Season
	if err := c.get(endpoint, params, &season); err != nil {
		return nil, err
	}
	return &season, nil
}

func (c *Client) GetEpisodeDetails(tvID int, seasonNumber int, episodeNumber int) (*Episode, error) {
	endpoint := fmt.Sprintf("%s/tv/%d/season/%d/episode/%d", tmdbAPIBase, tvID, seasonNumber, episodeNumber)
	params := url.Values{}
	params.Add("api_key", c.apiKey)

	var episode Episode
	if err := c.get(endpoint, params, &episode); err != nil {
		return nil, err
	}
	return &episode, nil
}

// Helper methods
func (c *Client) get(endpoint string, params url.Values, result interface{}) error {
	resp, err := c.client.Get(endpoint + "?" + params.Encode())
	if err != nil {
		return fmt.Errorf("TMDB API request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("TMDB API returned status: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("failed to decode TMDB response: %v", err)
	}

	return nil
}

func (c *Client) GetImageURL(path string) string {
	if path == "" {
		return ""
	}
	return tmdbImageBase + path
}
