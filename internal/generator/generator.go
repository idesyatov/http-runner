package generator

import (
    "fmt"
	"sync"
    "github.com/idesyatov/http-runner/pkg/httpclient"
)

type Generator struct {
    Client *httpclient.Client
}

func NewGenerator(client *httpclient.Client) *Generator {
    return &Generator{Client: client}
}

func (g *Generator) GenerateRequests(method, url string, count int) {
    var wg sync.WaitGroup
    for i := 0; i < count; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            resp, err := g.Client.SendRequest(method, url)
            if err != nil {
                fmt.Println("Error:", err)
                return
            }
            fmt.Println("Response Status:", resp.Status)
        }()
    }
    wg.Wait()
}
