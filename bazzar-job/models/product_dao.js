var constants = require('../constants');

var Sequelize = require('sequelize');
var sqldb = require('./mysql/DBModel');
var modelDB = sqldb.Product;

exp = {
    getById: function (id) {
        return modelDB.findOne({
            order: [
                ['id', 'DESC']],
            where: {
                id: id,
            }
        });
    },
}

module.exports = exp;