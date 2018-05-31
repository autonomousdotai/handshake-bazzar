var config = require('../../configs/index')
var Sequelize = require('sequelize');

var db = new Sequelize(
    config.mysql.database,
    config.mysql.username,
    config.mysql.password,
    config.mysql
);

var EthEvent = db.define('eth_event', {
    id: {
        type: Sequelize.INTEGER,
        autoIncrement: true,
        primaryKey: true,
        allowNull: false,
        unique: true,
    },
    address: Sequelize.STRING,
    event: Sequelize.STRING,
    value: Sequelize.STRING,
    block: Sequelize.INTEGER,
    log_index: Sequelize.INTEGER,
    date_created: Sequelize.DATE,
}, {
    tableName: 'eth_event',
    timestamps: false,
    underscored: true
});

var EthTx = db.define('eth_tx', {
    id: {
        type: Sequelize.INTEGER,
        autoIncrement: true,
        primaryKey: true,
        allowNull: false,
        unique: true,
    },
    user_id: Sequelize.INTEGER,
    hash: Sequelize.STRING,
    ref_type: Sequelize.STRING,
    ref_id: Sequelize.INTEGER,
    status: Sequelize.INTEGER,
    date_created: Sequelize.DATE,
    date_modified: Sequelize.DATE,
    from_address: Sequelize.STRING,
    to_address: Sequelize.STRING,
    tx_time: Sequelize.DATE,
    block: Sequelize.DECIMAL,
    value: Sequelize.DECIMAL,
    tx_fee: Sequelize.DECIMAL,
    raw_data: Sequelize.STRING,
}, {
    tableName: 'eth_tx',
    timestamps: false,
    underscored: true
});

var Product = db.define('product', {
    id: {
        type: Sequelize.INTEGER,
        autoIncrement: true,
        primaryKey: true,
        allowNull: false,
        unique: true,
    },
    date_created: Sequelize.DATE,
    date_modified: Sequelize.DATE,
    modified_user_id: Sequelize.INTEGER,
    created_user_id: Sequelize.INTEGER,
    hid: Sequelize.INTEGER,
    user_id: Sequelize.INTEGER,
    name: Sequelize.STRING,
    description: Sequelize.STRING,
    short_description: Sequelize.STRING,
    image: Sequelize.STRING,
    youtube_url: Sequelize.STRING,
    crowd_date: Sequelize.DATE,
    deliver_date: Sequelize.DATE,
    price: Sequelize.DECIMAL,
    goal: Sequelize.DECIMAL,
    balance: Sequelize.DECIMAL,
    shake_num: Sequelize.INTEGER,
    status: Sequelize.INTEGER
}, {
    tableName: 'crowd_funding',
    timestamps: false,
    underscored: true
});

var ProductShake = db.define('product_shake', {
    id: {
        type: Sequelize.INTEGER,
        autoIncrement: true,
        primaryKey: true,
        allowNull: false,
        unique: true,
    },
    date_created: Sequelize.DATE,
    date_modified: Sequelize.DATE,
    modified_user_id: Sequelize.INTEGER,
    created_user_id: Sequelize.INTEGER,
    user_id: Sequelize.INTEGER,
    crowd_funding_id: Sequelize.INTEGER,
    price: Sequelize.DECIMAL,
    quantity: Sequelize.INTEGER,
    amount: Sequelize.DECIMAL,
    status: Sequelize.INTEGER,
    address: Sequelize.STRING
}, {
    tableName: 'crowd_funding',
    timestamps: false,
    underscored: true
});


module.exports.db = db;
module.exports.EthEvent = EthEvent;
module.exports.EthTx = EthTx;
module.exports.Product = Product;
module.exports.ProductShake = ProductShake;
