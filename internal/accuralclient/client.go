package accuralclient

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/objects"
)

type AccuralClient struct {
	Cfg Config
}

func New(config Config) AccuralClient {
	return AccuralClient{Cfg: config}
}

func (c *AccuralClient) GetWithRetrues(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	for i := 0; i < c.Cfg.Retries && err != nil; i++ {
		time.Sleep(c.Cfg.Timeout)
		resp, err = http.Get(url)
	}
	return resp, err
}

func (c *AccuralClient) GetOrder(order string) (objects.AccuralOrder, error) {
	resp, err := c.GetWithRetrues("http//" + path.Join(c.Cfg.AccuralHost, "/api/orders/", order))
	if err != nil {
		return objects.AccuralOrder{}, errors.New("cant to get order from accural service")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error while read body: " + err.Error())
		return objects.AccuralOrder{}, errors.New("error while read body")
	}

	accuralOrder := objects.AccuralOrder{}
	if err := json.Unmarshal(body, &accuralOrder); err != nil {
		log.Println("error while unmarshal: " + err.Error())
		return objects.AccuralOrder{}, errors.New("error while unmarshal body")
	}

	return accuralOrder, nil
}
