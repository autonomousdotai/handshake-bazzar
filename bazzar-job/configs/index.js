module.exports = {
    mysql: {
        database: process.env.MYSQL_DATABASE,
        username: process.env.MYSQL_USERNAME,
        password: process.env.MYSQL_PASSWORD,
        host: process.env.MYSQL_HOST,
        dialect: 'mysql',
        timezone: process.env.MYSQL_TIMEZONE,
        pool: {
            max: 5,
            min: 0,
            idle: 10000
        }
    },
    timeAlive: 60,
    payableContractAddress: process.env.CONTRACT_ADDRESS,
    blockchainNetwork: process.env.CONTRACT_NETWORK,
}