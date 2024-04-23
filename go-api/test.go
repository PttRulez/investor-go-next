package main

import (
	"context"
	"fmt"

	"github.com/pttrulez/investor-go/internal/config"
	"github.com/pttrulez/investor-go/internal/repository/postgres"
)

func main() {
	cfg := config.MustLoad()
	repository, err := postgres.NewPostgresRepo(cfg.Pg)
	if err != nil {
		panic("Failed to initialize postgres repository: " + err.Error())
	}
	a, _ := repository.Moex.Shares.GetListByIds(context.Background(), []int{1, 2})
	fmt.Printf("SHARES:\n %+v", a)
}
