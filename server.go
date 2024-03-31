package main

import (
	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type Gear struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"not null"`
	Color string `gorm:"not null"`
	Size  string `gorm:"not null"`
	Price int    `gorm:"not null"`
}

func main() {
	db, err := gorm.Open(sqlite.Open("gear.db"))
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&Gear{})
	if err != nil {
		return
	}

	e := echo.New()

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db)
			return next(c)
		}
	})

	e.GET("/gears", getGears)
	e.GET("/gears/:id", getGear)
	e.PUT("/gears/:id", updateGear)
	e.POST("/gears", createGear)
	e.DELETE("/gears/:id", deleteGear)

	e.Logger.Fatal(e.Start(":8080"))
}

func getGears(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	var gears []Gear
	db.Find(&gears)
	return c.JSON(http.StatusOK, gears)
}

func getGear(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	id, _ := strconv.Atoi(c.Param("id"))
	var gear Gear
	db.First(&gear, id)
	return c.JSON(http.StatusOK, gear)
}

func updateGear(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	id, _ := strconv.Atoi(c.Param("id"))
	var gear Gear
	if err := db.First(&gear, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Gear not found"})
	}
	if err := c.Bind(&gear); err != nil {
		return err
	}
	db.Save(&gear)
	return c.JSON(http.StatusOK, gear)
}

func createGear(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	gear := new(Gear)
	if err := c.Bind(gear); err != nil {
		return err
	}
	db.Create(&gear)
	return c.JSON(http.StatusCreated, gear)
}

func deleteGear(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	id, _ := strconv.Atoi(c.Param("id"))
	var gear Gear
	if err := db.First(&gear, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Gear not found"})
	}
	db.Delete(&gear)
	return c.NoContent(http.StatusNoContent)
}
