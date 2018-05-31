var Sequelize = require('sequelize');

var sqldb = require('./mysql/DBModel');
var modelDB = sqldb.CrowdFundingShake;
var constants = require('../constants');

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
    updateHid: function (tx, id, hid, address) {
        address = address.toLowerCase()
        return modelDB.update(
            {
                hid: hid,
                address: address,
            },
            {
                where:
                    {
                        id: id
                    }
            }, {transaction: tx}
        );
    },
    updateStatus: function (tx, id, status) {
        address = address.toLowerCase()
        return modelDB.update(
            {
                status: status,
            },
            {
                where:
                    {
                        id: id
                    }
            }, {transaction: tx}
        );
    },
}

module.exports = exp;