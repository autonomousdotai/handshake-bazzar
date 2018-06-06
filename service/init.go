package service

import (
	"github.com/ninjadotorg/handshake-bazzar/dao"
	"github.com/ninjadotorg/handshake-bazzar/utils"
)

var fileUploadService = utils.GSService{}

// service
var productDao = dao.ProductDao{}
var productImageDao = dao.ProductImageDao{}
var productShakeDao = dao.ProductShakeDao{}
var ethTxDao = dao.EthTxDao{}

// template
var netUtil = utils.NetUtil{}
