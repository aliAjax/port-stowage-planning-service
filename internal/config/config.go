package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Addr          string
	SolveBudgetMS int
	LogJSON       bool
}

func Load() (Config, error) {
	c := Config{Addr: env("PORT_ADDR", ":8080"), SolveBudgetMS: 500, LogJSON: true}
	if v := os.Getenv("PORT_SOLVE_BUDGET_MS"); v != "" {
		n, e := strconv.Atoi(v)
		if e != nil || n < 10 {
			return c, fmt.Errorf("invalid PORT_SOLVE_BUDGET_MS")
		}
		c.SolveBudgetMS = n
	}
	return c, nil
}
func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
