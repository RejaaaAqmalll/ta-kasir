package produk

import (
	"fmt"
	"io"
	"net/http"
	"ta-kasir/base"
	"ta-kasir/config"
	"ta-kasir/helper"
	"ta-kasir/model"
	"ta-kasir/model/request"
	"ta-kasir/model/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddProduk(c *gin.Context) {
	dataJWT, err := helper.GetClaims(c)
	
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
			Status:  http.StatusUnauthorized,
			Error:   err,
			Message: base.NoUserLogin,
			Data:    nil,
		})
	}

	formAddProduk := request.AddProduk{}

	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	err = c.ShouldBind(&formAddProduk)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: base.EmpetyField,
			Data:    nil,
		})
		return
	}

	// validasi input file harus berupa gambar
	src, err := file.Open()
	if err != nil{
		// log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	defer src.Close()

	buffer := make([]byte, 261)
	_, err = src.Read(buffer)

	if err != nil && err != io.EOF {
		// log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	
	// get mime type
	kind := http.DetectContentType(buffer)
	if kind == "" || !helper.IsSupportedImageFormat(kind) {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: base.FileNotSupported,
			Data:    nil,
		})
		return
	}

	fileName := helper.GenerateFilename(file.Filename)

	err = helper.SaveFile(src, fileName)
	if err != nil {
		// log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
			Status:  http.StatusInternalServerError,
			Error:   err,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	link := fmt.Sprintf("storage/%s", fileName)

	db := config.ConnectDatabase()

	err = db.Debug().Where("email = ?", dataJWT.Email).
	Where("role = ?", 1).Where("hapus = ?", 0).First(&model.User{}).Error

	if  err != nil && err == gorm.ErrRecordNotFound {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
			Status:  http.StatusUnauthorized,
			Error:   err,
			Message: base.ShouldAdmin,
			Data:    nil,
		})
		return
	}

	var produk  = model.Produk{
	NamaProduk: formAddProduk.NamaProduk,
	Harga:      formAddProduk.Harga,
	Stok:       formAddProduk.Stok,
	Gambar: 	link,	
	}

	err = db.Debug().Create(&produk).Error

	if err != nil {
		// log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
			Status:  http.StatusInternalServerError,
			Error:   err,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	// fmt.Println(link)
	finalLink := "http://127.0.0.1:8080/" + link
	// fmt.Println(finalLink)
	c.JSON(http.StatusOK, response.Response{
		Status:  http.StatusOK,
		Error:   nil,
		Message: base.SuccessAddProduk,
		Data:    gin.H{
			"data_produk": produk,
			"link":        finalLink,
		},
	})
}

func EditProduk(c *gin.Context)  {
	dataJWT, err := helper.GetClaims(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
			Status:  http.StatusUnauthorized,
			Error:   err,
			Message: base.NoUserLogin,
			Data:    nil,
		})
		return
	}

	idProduk := c.Param("id")
	if idProduk == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status: http.StatusBadRequest,
			Error: nil,
			Message: base.ParamEmpty,
			Data: nil,
		})
		return
	}

	formEditProduk := request.EditProduk{}
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.ShouldBind(&formEditProduk)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: base.EmpetyField,
			Data:    nil,
		})
		return
	}

	// validasi input file harus berupa gambar
	src, err := file.Open()
	if err != nil{
		// log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	defer src.Close()

	buffer := make([]byte, 261)
	_, err = src.Read(buffer)

	if err != nil && err != io.EOF {
		// log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	
	// get mime type
	kind := http.DetectContentType(buffer)
	if kind == "" || !helper.IsSupportedImageFormat(kind) {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status:  http.StatusBadRequest,
			Error:   err,
			Message: base.FileNotSupported,
			Data:    nil,
		})
		return
	}

	fileName := helper.GenerateFilename(file.Filename)

	err = helper.SaveFile(src, fileName)
	if err != nil {
		// log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
			Status:  http.StatusInternalServerError,
			Error:   err,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	link := fmt.Sprintf("storage/%s", fileName)

	db := config.ConnectDatabase()

	err = db.Debug().Where("email = ?", dataJWT.Email).
	Where("role = ?", 1).Where("hapus = ?", 0).First(&model.User{}).Error

	if  err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
			Status:  http.StatusUnauthorized,
			Error:   err,
			Message: base.ShouldAdmin,
			Data:    nil,
		})
		return
	}

	var produk  = model.Produk{
		NamaProduk: formEditProduk.NamaProduk,
		Harga: formEditProduk.Harga,
		Stok: formEditProduk.Stok,
		Gambar: link,
	}

	err = db.Debug().Model(model.Produk{}).
	Where("id_produk = ?", idProduk).
	Updates(&produk).Error

	if err != nil {
		// log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
			Status:  http.StatusInternalServerError,
			Error:   err,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	finalLink := "http://127.0.0.1:8080/" + link
	// fmt.Println(finalLink)
	c.JSON(http.StatusOK, response.Response{
		Status:  http.StatusOK,
		Error:   nil,
		Message: base.SuccessEditPorduk,
		Data:    gin.H{
			"data_produk": produk,
			"link":        finalLink,
		},
	})
}

func DeleteProduk(c *gin.Context)  {
	idProduk := c.Param("id")

	if idProduk == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Status: http.StatusBadRequest,
			Error: nil,
			Message: base.ParamEmpty,
			Data: nil,
		})
		return
	}

	dataJWT, err := helper.GetClaims(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
			Status: http.StatusUnauthorized,
			Error: err,
			Message: base.NoUserLogin,
			Data: nil,
		})
		return
	}

	db := config.ConnectDatabase()
	err = db.Where("email = ?", dataJWT.Email).
	Where("role = ?", 1).Where("hapus = ?", 0).
	First(&model.User{}).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
			Status: http.StatusUnauthorized,
			Error: err,
			Message: base.ShouldAdmin,
			Data: nil,
		})
		return
	}

	err = db.Debug().Model(model.Produk{}).
	Where("id_produk = ?", idProduk).Update("hapus", 1).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
			Status: http.StatusInternalServerError,
			Error: err,
			Message: err.Error(),
			Data: nil,
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Status: http.StatusOK,
		Error:  nil,
		Message: base.SuccessDeleteProduk,
		Data: nil,
	})
}