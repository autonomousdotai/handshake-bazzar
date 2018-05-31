var cron = require('node-cron');

var ethEventDAO = require('./models/eth_event_dao');
var productDao = require('./models/product_dao');
var productShakeDao = require('./models/product_shake_dao');
var ethTxDAO = require('./models/eth_tx_dao');


var sqldb = require('./models/mysql/DBModel');
var modelDB = sqldb.db;

var config = require('./configs/index')
var constants = require('./constants');


var Web3 = require('web3');
web3 = new Web3(new Web3.providers.HttpProvider(config.blockchainNetwork));

var payableContractJson = require('./contracts/PayableHandshake.json');
var payableContractAddress = config.payableContractAddress
var payableContractEventNames = ['__init', '__shake', '__deliver', '__cancel', '__reject', '__accept', '__withdraw'];
var payableContractInstance = new web3.eth.Contract(payableContractJson.abi, payableContractAddress);

console.log('Events by blockchainNetwork: ' + config.blockchainNetwork);
console.log('Events by payableContractAddress: ' + payableContractAddress);

function parseOffchain(offchain) {
    let values = offchain.replace(/\u0000/g, '').split("_")
    console.log(values)
    if (values.length >= 2) {
        return [values[0].trim(), values[1].trim()];
    } else {
        return null;
    }
}

async function processEventObj(contractAddress, eventName, eventObj) {
    let tx = await modelDB.transaction();
    try {
        console.log("processEventObj", contractAddress, eventName, eventObj);

        let tx_hash = eventObj.transactionHash.toLowerCase()
        let txr = await web3.eth.getTransactionReceipt(tx_hash);

        await ethEventDAO.create(tx, contractAddress, eventName, JSON.stringify(eventObj), eventObj.blockNumber, eventObj.logIndex);

        switch (contractAddress) {
            case payableContractAddress: {
                switch (eventName) {
                    case '__init': {
                        console.log("__init hid = " + eventObj.returnValues.hid);
                        console.log("__init offchain = " + eventObj.returnValues.offchain);
                        const hid = eventObj.returnValues.hid;
                        const offchain = eventObj.returnValues.offchain;
                        if (hid == undefined || offchain == undefined) {
                            console.log("__init missing parameters");
                            break;
                        }
                        let offchainStr = Web3.utils.toAscii(offchain);
                        let offchains = parseOffchain(offchainStr);
                        console.log("__init offchains", offchains);
                        let offchainType = offchains[0];
                        if (offchainType == constants.OFFCHAIN_TYPE_SHAKE) {
                            let offchainId = parseInt(offchains[1]);
                            let productShake = await productShakeDao.getById(offchainId);
                            if (productShake == null) {
                                console.log("__init productDao.getById NULL", offchainId);
                                break;
                            }
                            let tns = await ethTxDAO.getByHash(tx_hash);
                            if (tns == null) {
                                tns = await ethTxDAO.create(productShake.user_id, tx_hash, 'payable_init', productShake.id);
                            }
                        }
                    }
                        break;
                    case '__shake': {
                        console.log("__shake hid = " + eventObj.returnValues.hid);
                        console.log("__shake offchain = " + eventObj.returnValues.offchain);
                        const hid = eventObj.returnValues.hid;
                        const offchain = eventObj.returnValues.offchain;
                        if (hid == undefined || offchain == undefined) {
                            console.log("__shake missing parameters");
                            break;
                        }
                        let offchainStr = Web3.utils.toAscii(offchain);
                        let offchains = parseOffchain(offchainStr);
                        console.log("__shake offchains", offchains);
                        let offchainType = offchains[0];
                        if (offchainType == constants.OFFCHAIN_TYPE_SHAKE) {
                            let productShakeId = parseInt(offchains[1]);
                            let productShake = await productShakeDao.getById(productShakeId);
                            if (productShake == null) {
                                console.log("__shake productShakeDAO.getById NULL", productShakeId);
                                break;
                            }
                            console.log("__shake productShakeDAO.getById OK", productShakeId);

                            await productShakeDao.updateHid(tx, productShake.id, hid, txr.address)
                            console.log("__shake productShakeDAO.updateHid OK", productShake.id, hid, txr.address)

                            await productShakeDao.updateStatus(tx, productShake.id, constants.ORDER_STATUS_SHAKED)
                            console.log("__shake productShakeDAO.updateStatus OK", productShake.id, constants.ORDER_STATUS_SHAKED)

                            let tns = await ethTxDAO.getByHash(tx_hash);
                            if (tns == null) {
                                tns = await ethTxDAO.create(productShake.user_id, tx_hash, 'payable_shake', productShake.id);
                            }
                        }
                    }
                        break;
                    case '__deliver': {
                        console.log("__shake hid = " + eventObj.returnValues.hid);
                        console.log("__shake offchain = " + eventObj.returnValues.offchain);
                        const hid = eventObj.returnValues.hid;
                        const offchain = eventObj.returnValues.offchain;
                        if (hid == undefined || offchain == undefined) {
                            console.log("__shake missing parameters");
                            break;
                        }
                        let offchainStr = Web3.utils.toAscii(offchain);
                        let offchains = parseOffchain(offchainStr);
                        console.log("__shake offchains", offchains);
                        let offchainType = offchains[0];
                        if (offchainType == constants.OFFCHAIN_TYPE_SHAKE) {
                            let productShakeId = parseInt(offchains[1]);
                            let productShake = await productShakeDao.getById(productShakeId);
                            if (productShake == null) {
                                console.log("__shake productShakeDAO.getById NULL", productShakeId);
                                break;
                            }
                            console.log("__shake productShakeDAO.getById OK", productShakeId);

                            await productShakeDao.updateStatus(tx, productShake.id, constants.ORDER_STATUS_DELIVERED)
                            console.log("__shake productShakeDAO.updateStatus OK", productShake.id, constants.ORDER_STATUS_DELIVERED)

                            let tns = await ethTxDAO.getByHash(tx_hash);
                            if (tns == null) {
                                tns = await ethTxDAO.create(productShake.user_id, tx_hash, 'payable_shake', productShake.id);
                            }
                        }
                    }
                        break;
                    case '__cancel': {
                        console.log("__shake hid = " + eventObj.returnValues.hid);
                        console.log("__shake offchain = " + eventObj.returnValues.offchain);
                        const hid = eventObj.returnValues.hid;
                        const offchain = eventObj.returnValues.offchain;
                        if (hid == undefined || offchain == undefined) {
                            console.log("__shake missing parameters");
                            break;
                        }
                        let offchainStr = Web3.utils.toAscii(offchain);
                        let offchains = parseOffchain(offchainStr);
                        console.log("__shake offchains", offchains);
                        let offchainType = offchains[0];
                        if (offchainType == constants.OFFCHAIN_TYPE_SHAKE) {
                            let productShakeId = parseInt(offchains[1]);
                            let productShake = await productShakeDao.getById(productShakeId);
                            if (productShake == null) {
                                console.log("__shake productShakeDAO.getById NULL", productShakeId);
                                break;
                            }
                            console.log("__shake productShakeDAO.getById OK", productShakeId);

                            await productShakeDao.updateStatus(tx, productShake.id, constants.ORDER_STATUS_CANCELED)
                            console.log("__shake productShakeDAO.updateStatus OK", productShake.id, constants.ORDER_STATUS_CANCELED)

                            let tns = await ethTxDAO.getByHash(tx_hash);
                            if (tns == null) {
                                tns = await ethTxDAO.create(productShake.user_id, tx_hash, 'payable_cancel', productShake.id);
                            }
                        }
                    }
                        break;
                    case '__reject': {
                        console.log("__shake hid = " + eventObj.returnValues.hid);
                        console.log("__shake offchain = " + eventObj.returnValues.offchain);
                        const hid = eventObj.returnValues.hid;
                        const offchain = eventObj.returnValues.offchain;
                        if (hid == undefined || offchain == undefined) {
                            console.log("__shake missing parameters");
                            break;
                        }
                        let offchainStr = Web3.utils.toAscii(offchain);
                        let offchains = parseOffchain(offchainStr);
                        console.log("__shake offchains", offchains);
                        let offchainType = offchains[0];
                        if (offchainType == constants.OFFCHAIN_TYPE_SHAKE) {
                            let productShakeId = parseInt(offchains[1]);
                            let productShake = await productShakeDao.getById(productShakeId);
                            if (productShake == null) {
                                console.log("__shake productShakeDAO.getById NULL", productShakeId);
                                break;
                            }
                            console.log("__shake productShakeDAO.getById OK", productShakeId);

                            await productShakeDao.updateStatus(tx, productShake.id, constants.ORDER_STATUS_REJECTED)
                            console.log("__shake productShakeDAO.updateStatus OK", productShake.id, constants.ORDER_STATUS_REJECTED)

                            let tns = await ethTxDAO.getByHash(tx_hash);
                            if (tns == null) {
                                tns = await ethTxDAO.create(productShake.user_id, tx_hash, 'payable_reject', productShake.id);
                            }
                        }
                    }
                        break;
                    case '__accept': {
                        console.log("__shake hid = " + eventObj.returnValues.hid);
                        console.log("__shake offchain = " + eventObj.returnValues.offchain);
                        const hid = eventObj.returnValues.hid;
                        const offchain = eventObj.returnValues.offchain;
                        if (hid == undefined || offchain == undefined) {
                            console.log("__shake missing parameters");
                            break;
                        }
                        let offchainStr = Web3.utils.toAscii(offchain);
                        let offchains = parseOffchain(offchainStr);
                        console.log("__shake offchains", offchains);
                        let offchainType = offchains[0];
                        if (offchainType == constants.OFFCHAIN_TYPE_SHAKE) {
                            let productShakeId = parseInt(offchains[1]);
                            let productShake = await productShakeDao.getById(productShakeId);
                            if (productShake == null) {
                                console.log("__shake productShakeDAO.getById NULL", productShakeId);
                                break;
                            }
                            console.log("__shake productShakeDAO.getById OK", productShakeId);

                            await productShakeDao.updateStatus(tx, productShake.id, constants.ORDER_STATUS_ACCEPTED)
                            console.log("__shake productShakeDAO.updateStatus OK", productShake.id, constants.ORDER_STATUS_ACCEPTED)

                            let tns = await ethTxDAO.getByHash(tx_hash);
                            if (tns == null) {
                                tns = await ethTxDAO.create(productShake.user_id, tx_hash, 'payable_accept', productShake.id);
                            }
                        }
                    }
                        break;
                    case '__withdraw': {
                        console.log("__shake hid = " + eventObj.returnValues.hid);
                        console.log("__shake offchain = " + eventObj.returnValues.offchain);
                        const hid = eventObj.returnValues.hid;
                        const offchain = eventObj.returnValues.offchain;
                        if (hid == undefined || offchain == undefined) {
                            console.log("__shake missing parameters");
                            break;
                        }
                        let offchainStr = Web3.utils.toAscii(offchain);
                        let offchains = parseOffchain(offchainStr);
                        console.log("__shake offchains", offchains);
                        let offchainType = offchains[0];
                        if (offchainType == constants.OFFCHAIN_TYPE_SHAKE) {
                            let productShakeId = parseInt(offchains[1]);
                            let productShake = await productShakeDao.getById(productShakeId);
                            if (productShake == null) {
                                console.log("__shake productShakeDAO.getById NULL", productShakeId);
                                break;
                            }
                            console.log("__shake productShakeDAO.getById OK", productShakeId);

                            await productShakeDao.updateStatus(tx, productShake.id, constants.ORDER_STATUS_WITHDRAWED)
                            console.log("__shake productShakeDAO.updateStatus OK", productShake.id, constants.ORDER_STATUS_WITHDRAWED)

                            let tns = await ethTxDAO.getByHash(tx_hash);
                            if (tns == null) {
                                tns = await ethTxDAO.create(productShake.user_id, tx_hash, 'payable_withdraw', productShake.id);
                            }
                        }
                    }
                        break;
                }
            }
                break;
        }
        tx.commit();
    } catch (err) {
        console.log('processEventObj', err);
        tx.rollback();
    }
}

function asyncGetPastEvents(contract, contractAddress, eventName, fromBlock) {
    return new Promise(function (resolve, reject) {
        contract.getPastEvents(eventName, {
            filter: {_from: contractAddress},
            fromBlock: fromBlock,
            toBlock: 'latest'

        }, function (error, events) {
            console.log(eventName + " getPastEvents OK")
            if (error != null) {
                reject(error);
            } else {
                resolve(events);
            }
        });
    })
}

async function asyncScanEventLog(contract, contractAddress, eventName) {
    let lastEventLog = await ethEventDAO.getLastLogByName(contractAddress, eventName);
    var fromBlock = 0;
    if (lastEventLog != null) {
        fromBlock = lastEventLog.block + 1;
    }
    console.log(eventName + " fromBlock = " + fromBlock);
    let events = await asyncGetPastEvents(contract, contractAddress, eventName, fromBlock);
    for (var i = 0; i < events.length; i++) {
        const eventObj = events[i];
        console.log(eventObj);
        let checkEventLog = await ethEventDAO.getByBlock(contractAddress, eventObj.blockNumber, eventObj.logIndex);
        if (checkEventLog == null) {
            await processEventObj(contractAddress, eventName, eventObj);
        }
    }

}

async function processTx(id, user_id, hash, ref_type, ref_id, date_created) {
    let tx = await modelDB.transaction();
    try {
        let txr = null;
        try {
            txr = await web3.eth.getTransactionReceipt(hash);
        } catch (err) {
            console.log('error', err)
            txr = null;
        }
        let is_failed = false
        if (txr == null) {
            let now = new Date()
            if (now - date_created > 24 * 60 * 60 * 1000) {
                is_failed = true
            } else {
                await ethTxDAO.updateStatus(tx, hash, 0);
                console.log('txr is null', hash);
                tx.commit();
                return
            }
        } else {
            console.log('txr is ok', txr);
            let txrJson = JSON.stringify(txr);
            await ethTxDAO.updateInfo(tx, id, txr.from, txr.to, new Date(), txr.blockNumber, 0, 0, txrJson);
            is_failed = (txr.status == '1' || txr.status == '0x1') ? false : true;
        }
        if (is_failed) {
            if (txr != null) {
                await ethTxDAO.updateStatus(tx, hash, 2);
            } else {
                await ethTxDAO.updateStatus(tx, hash, 3);
            }
        } else {
            await ethTxDAO.updateStatus(tx, hash, 1);
        }
        switch (ref_type) {
            case 'payable_shake':{
                if (is_failed) {
                    let productShake = await productShakeDao.getById(ref_id);
                    if (productShake == null) {
                        console.log(ref_type + ' productDao.getById NULL', ref_id);
                        break;
                    }
                    console.log(ref_type + ' productDao.getById OK', ref_id);
                    await productShakeDao.updateStatus(tx, productShake.id, constants.ORDER_STATUS_SHAKED_FAILED);
                    console.log(ref_type + ' productDao.updateStatus OK', productShake.id, constants.ORDER_STATUS_SHAKED_FAILED);
                }
            }
                break;
            case 'payable_deliver':{
                if (is_failed) {
                    let productShake = await productShakeDao.getById(ref_id);
                    if (productShake == null) {
                        console.log(ref_type + ' productDao.getById NULL', ref_id);
                        break;
                    }
                    console.log(ref_type + ' productDao.getById OK', ref_id);
                    await productShakeDao.updateStatus(tx, productShake.id, constants.ORDER_STATUS_DELIVERED_FAILED);
                    console.log(ref_type + ' productDao.updateStatus OK', productShake.id, constants.ORDER_STATUS_DELIVERED_FAILED);
                }
            }
                break;
            case 'payable_cancel':{
                if (is_failed) {
                    let productShake = await productShakeDao.getById(ref_id);
                    if (productShake == null) {
                        console.log(ref_type + ' productDao.getById NULL', ref_id);
                        break;
                    }
                    console.log(ref_type + ' productDao.getById OK', ref_id);
                    await productShakeDao.updateStatus(tx, productShake.id, constants.ORDER_STATUS_CANCELED_FAILED);
                    console.log(ref_type + ' productDao.updateStatus OK', productShake.id, constants.ORDER_STATUS_CANCELED_FAILED);
                }
            }
                break;
            case 'payable_reject':{
                if (is_failed) {
                    let productShake = await productShakeDao.getById(ref_id);
                    if (productShake == null) {
                        console.log(ref_type + ' productDao.getById NULL', ref_id);
                        break;
                    }
                    console.log(ref_type + ' productDao.getById OK', ref_id);
                    await productShakeDao.updateStatus(tx, productShake.id, constants.ORDER_STATUS_REJECTED_FAILED);
                    console.log(ref_type + ' productDao.updateStatus OK', productShake.id, constants.ORDER_STATUS_REJECTED_FAILED);
                }
            }
                break;
            case 'payable_accept':{
                if (is_failed) {
                    let productShake = await productShakeDao.getById(ref_id);
                    if (productShake == null) {
                        console.log(ref_type + ' productDao.getById NULL', ref_id);
                        break;
                    }
                    console.log(ref_type + ' productDao.getById OK', ref_id);
                    await productShakeDao.updateStatus(tx, productShake.id, constants.ORDER_STATUS_ACCEPTED_FAILED);
                    console.log(ref_type + ' productDao.updateStatus OK', productShake.id, constants.ORDER_STATUS_ACCEPTED_FAILED);
                }
            }
                break;
            case 'payable_withdraw':{
                if (is_failed) {
                    let productShake = await productShakeDao.getById(ref_id);
                    if (productShake == null) {
                        console.log(ref_type + ' productDao.getById NULL', ref_id);
                        break;
                    }
                    console.log(ref_type + ' productDao.getById OK', ref_id);
                    await productShakeDao.updateStatus(tx, productShake.id, constants.ORDER_STATUS_WITHDRAWED_FAILED);
                    console.log(ref_type + ' productDao.updateStatus OK', productShake.id, constants.ORDER_STATUS_WITHDRAWED_FAILED);
                }
            }
                break;
        }
        tx.commit();
    } catch (err) {
        console.log('error', err)
        tx.rollback();
    }
}

async function cronJob() {
    console.log('running a task every minute at ' + new Date());
    console.log('process ether tx');
    let results = await ethTxDAO.getListUnTx();
    for (var i = 0; i < results.length; i++) {
        var result = results[i];
        await processTx(result.id, result.user_id, result.hash, result.ref_type, result.ref_id, result.date_created);
    }
    console.log('process ether events');
    if (payableContractAddress != '') {
        for (var i = 0; i < payableContractEventNames.length; i++) {
            var eventName = payableContractEventNames[i];
            await asyncScanEventLog(payableContractInstance, payableContractAddress, eventName);
        }
    }
}

cronJob();

cron.schedule('* * * * *', async function () {
    await cronJob();
});

