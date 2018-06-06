package service

import (
	"github.com/autonomousdotai/handshake-bazzar/bazzar-service/dao"
	"github.com/autonomousdotai/handshake-bazzar/bazzar-service/utils"
)

var fileUploadService = utils.GSService{}
// service
var productDao = dao.ProductDao{}
var productImageDao = dao.ProductImageDao{}
var productShakeDao = dao.ProductShakeDao{}
var ethTxDao = dao.EthTxDao{}
// template
var netUtil = utils.NetUtil{}