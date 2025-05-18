package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/shirou/gopsutil/v3/cpu"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type CPUUsage struct {
	CPU       string  `json:"cpu"`
	Usage     float64 `json:"usage"`
	Timestamp int64   `json:"timestamp"`
}

func getCPUUsage() ([]CPUUsage, error) {
	// Obtenir le nombre de cœurs logiques
	counts, err := cpu.Counts(true)
	if err != nil {
		return nil, fmt.Errorf("impossible d'obtenir le nombre de cœurs: %v", err)
	}

	// Obtenir les pourcentages d'utilisation pour chaque cœur
	percent, err := cpu.Percent(500*time.Millisecond, true)
	if err != nil {
		return nil, fmt.Errorf("impossible d'obtenir l'utilisation du CPU: %v", err)
	}

	// Préparer les données de sortie
	usage := make([]CPUUsage, counts)
	now := time.Now().Unix()

	for i := 0; i < counts; i++ {
		cpuName := fmt.Sprintf("CPU%d", i)
		if i == 0 && counts == 1 {
			cpuName = "CPU"
		}
		usage[i] = CPUUsage{
			CPU:       cpuName,
			Usage:     percent[i],
			Timestamp: now,
		}
	}

	return usage, nil
}

// Handler pour l'utilisation du CPU
func getCPUUsageHandler(c echo.Context) error {
	usage, err := getCPUUsage()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	// Renvoyer un objet avec une propriété cpus pour correspondre à ce que le frontend attend
	return c.JSON(http.StatusOK, map[string]interface{}{
		"cpus": usage,
	})
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Templates
	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = t

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", nil)
	})

	e.GET("/cpu-usage", getCPUUsageHandler)

	// Démarrer le serveur
	e.Logger.Fatal(e.Start(":1323"))
}
