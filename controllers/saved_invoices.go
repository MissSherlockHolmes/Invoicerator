package controllers

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func ShowSavedInvoicesPage(c *gin.Context) {
    c.HTML(http.StatusOK, "saved_invoices.html", nil)
}
