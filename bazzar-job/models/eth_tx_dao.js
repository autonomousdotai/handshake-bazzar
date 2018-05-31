var sqldb = require('./mysql/DBModel');
var modelDB = sqldb.EthTx;

exp = {
    create: function (customer_id, hash, ref_type, ref_id, from_address, to_address, tx_time, block, value, tx_fee, raw_data) {
        hash = hash.toLowerCase()
        return modelDB.create({
            customer_id: customer_id,
            hash: hash,
            ref_type: ref_type,
            ref_id: ref_id,
            status: 1,
            date_created: new Date(),
            date_modified: new Date(),
            from_address: from_address,
            to_address: to_address,
            tx_time: tx_time,
            block: block,
            raw_data: raw_data,
        });
    },
    getById: function (id) {
        return modelDB.findOne({
            order: [
                ['id', 'DESC']],
            where: {
                id: id,
            }
        });
    },
    getByHash: function (hash) {
        hash = tx_hash.toLowerCase()
        return modelDB.findOne({
            order: [
                ['id', 'DESC']],
            where: {
                hash: hash,
            }
        });
    },
    getListUnTx: function () {
        return modelDB.findAll({
            order: [
                ['date_modified', 'ASC']],
            where: {
                status: 0,
            }
        });
    },
    updateStatus: function (tx, hash, status) {
        hash = hash.toLowerCase()
        return modelDB.update(
            {
                status: status,
                date_modified: new Date(),
            },
            {
                where:
                    {
                        hash: hash,
                        status: 0,
                    }
            }, {transaction: tx}
        );
    },
    updateInfo: function (tx, id, from_address, to_address, tx_time, block, value, tx_fee, raw_data) {
        from_address = from_address.toLowerCase()
        to_address = to_address.toLowerCase()
        return modelDB.update(
            {
                from_address: from_address,
                to_address: to_address,
                tx_time: tx_time,
                block: block,
                raw_data: raw_data,
            },
            {
                where:
                    {
                        id: id,
                    }
            }, {transaction: tx}
        );
    }

}

module.exports = exp;