package providers

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
)

type CommonProvider struct {
	*Config
	apiURL string
}

type CommonResult struct {
	URL   string `json:"url"`
	Error string `json:"error"`
}

type CommonPaginationResult struct {
	Blocks   uint `json:"blocks"`
	PageSize uint `json:"pageSize"`
	Pages    uint `json:"pages"`
}

type CommonAPIResult []struct {
	API string `json:"cdx-api"`
}

func NewCommonProvider(config *Config) (Provider, error) {
	c := CommonProvider{Config: config}

	// Fetch the list of available CommonCrawl Api URLs.
	resp, err := c.MakeRequest("http://index.commoncrawl.org/collinfo.json")
	if err != nil {
		if c.Config.Verbose {
			fmt.Printf("[-] Error on request collinfo.json: %s\n", err)
		}
		return nil, err
	}

	defer resp.Body.Close()

	var apiResult CommonAPIResult;
	err = json.NewDecoder(resp.Body).Decode(&apiResult);
	if  err != nil || len(apiResult) < 1{
		if c.Config.Verbose {
			fmt.Printf("[-] Erorr on response collinfo.json: %s\n", err)
		}
		err = errors.New("[-] Commoncrawl, erorr on response collinfo.json.")
		return nil, err
	}

	c.apiURL = apiResult[0].API
	return &c, nil
}

func (c *CommonProvider) formatURL(domain string, page uint) string {
	if c.IncludeSubdomains {
		domain = "*." + domain
	}

	return fmt.Sprintf("%s?url=%s/*&output=json&fl=url&page=%d", c.apiURL, domain, page)
}

// Fetch the number of pages.
func (c *CommonProvider) getPagination(domain string) (*CommonPaginationResult, error) {
	url := fmt.Sprintf("%s&showNumPages=true", c.formatURL(domain, 0))

	resp, err := c.MakeRequest(url)
	if err != nil {
		if c.Config.Verbose {
			fmt.Printf("[-] Error request pagination: %s\n", err)
		}
		return nil, err
	}

	defer resp.Body.Close()

	var paginationResult CommonPaginationResult
	if err = json.NewDecoder(resp.Body).Decode(&paginationResult); err != nil {
		if c.Config.Verbose {
			fmt.Printf("[-] Error response pagination: %s\n", err)
		}
		return nil, err
	}

	return &paginationResult, nil
}

func (c *CommonProvider) Fetch(domain string, results chan<- string) error {
	pagination, err := c.getPagination(domain)
	if err != nil {
		if c.Config.Verbose {
			fmt.Printf("[-] failed to fetch common pagination: %s\n", err)
		}
		return fmt.Errorf("failed to fetch common pagination: %s", err)
	}

	for page := uint(0); page < pagination.Pages; page++ {
		resp, err := c.MakeRequest(c.formatURL(domain, page))
		if err != nil {
			if c.Config.Verbose {
				fmt.Printf("[-] failed to fetch common results page %d: %s\n", page, err)
			}
			return fmt.Errorf("failed to fetch common results page %d: %s", page, err)
		}

		sc := bufio.NewScanner(resp.Body)
		for sc.Scan() {
			var result CommonResult
			if err := json.Unmarshal(sc.Bytes(), &result); err != nil {
				_ = resp.Body.Close()
				if c.Config.Verbose {
					fmt.Printf("[-] failed to decode common results for page %d: %s\n", page, err)
				}
				return fmt.Errorf("failed to decode common results for page %d: %s", page, err)
			}

			if result.Error != "" {
				if c.Config.Verbose {
					fmt.Printf("[-] received an error from common api: %s\n", result.Error)
				}
				return fmt.Errorf("received an error from common api: %s", result.Error)
			}

			if c.Config.Verbose {
				fmt.Printf("[Common] domain %v : %v\n", domain, len(result.URL))
			}
			results <- result.URL
		}

		_ = resp.Body.Close()
	}

	return nil
}
