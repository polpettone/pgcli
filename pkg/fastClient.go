package pkg

import (
	"net/http"
	"sort"
)

type result struct {
	index int
	res http.Response
	err error
}


func BoundedParallelRequestCall(client http.Client, requests []http.Request, concurrencyLimit int) []result {
	semaphoreChan := make(chan struct{}, concurrencyLimit)
	resultsChan := make(chan *result)

	defer func() {
		close(semaphoreChan)
		close(resultsChan)
	}()

	for i, request := range requests {
		go func (i int, request http.Request) {
			semaphoreChan <- struct{}{}
			res, err := client.Do(&request)
			result := &result{i, *res, err}
			resultsChan <- result
			<-semaphoreChan
		} (i,request)
	}

	var results []result

	for {
		result := <-resultsChan
		results = append(results, *result)

		if len(results) == len(requests) {
			break
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].index < results[j].index
	})
	return results
}