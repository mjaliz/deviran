package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/mjaliz/deviran/internal/initializers"
	"github.com/mjaliz/deviran/internal/message"
	"github.com/mjaliz/deviran/internal/models"
	"net/http"
)

func GetCurrencies(c echo.Context) error {
	var currencies []models.Currency
	if err := initializers.DB.Find(&currencies).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, message.StatusInternalServerErrorMessage())
	}
	return c.JSON(http.StatusOK, message.StatusOkMessage(currencies, ""))
}

func GetCurrency(c echo.Context) error {
	currencyId := c.Param("currency_id")
	var currencyPrices []models.CurrencyPrice

	if err := initializers.DB.Where("currency_id = ?", currencyId).Find(&currencyPrices).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, message.StatusInternalServerErrorMessage())
	}
	return c.JSON(http.StatusOK, message.StatusOkMessage(currencyPrices, ""))

}
