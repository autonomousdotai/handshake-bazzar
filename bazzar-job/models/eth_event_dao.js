var sqldb = require('./mysql/DBModel');
var model = sqldb.EthEvent;

exp = {
    getByBlock: function (address, block, log_index) {
        address = address.toLowerCase()
        return model
            .findOne({
                order: [
                    ['id', 'DESC']],
                where: {
                    address: address,
                    block: block,
                    log_index: log_index,
                }
            });
    },
    getLastLogByName: function (address, event) {
        address = address.toLowerCase()
        return model
            .findOne({
                order: [
                    ['id', 'DESC']],
                where: {
                    address: address,
                    event: event,
                }
            });
    },
    create: function (tx, address, event, value, block, log_index) {
        address = address.toLowerCase()
        return model
            .create(
                {
                    address: address,
                    event: event,
                    value: value,
                    block: block,
                    log_index: log_index,
                    date_created: new Date(),
                }, {transaction: tx}
            )
    }

}

module.exports = exp;