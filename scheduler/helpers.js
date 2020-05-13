const fetch = require("node-fetch");
const grpc = require('grpc');
const protoLoader = require('@grpc/proto-loader');

module.exports = {
    newPricer,
    graphQuery,
    grpcStreamQuery,
    newPriceRequest
}

function newPricer() {
    // grpc client
    var PROTO_PATH = __dirname + '/../gobs/api/protos/pricer.proto';
    var packageDefinition = protoLoader.loadSync(PROTO_PATH, {
        keepCase: true,
        longs: String,
        enums: String,
        defaults: true,
        oneofs: true
    });
    var gobs = grpc.loadPackageDefinition(packageDefinition).gobs;
    var client = new gobs.PricerService(':5050', grpc.credentials.createInsecure());
    return client;
}

function graphQuery(url, queryString, callback, errHandler) {
  fetch(url, {
    "method": "POST",
    "headers": {
      "content-type": "application/json"
    },
    "body": JSON.stringify({query: queryString})
  })
  .then(res => res.json())
  .then(jsonres => {
    // console.json(JSON.stringify(jsonres));
    if(callback && typeof callback === "function") {
      callback(jsonres);
    }
  })
  .catch(err => {
    console.json(JSON.stringify(err));
    if(errHandler && typeof errHandler === "function") {
      errHandler(err);
    }
  });
}

function grpcStreamQuery(stub, requests, callback, errHandler, endHandler) {
    stub.on('data', function(response) {
        console.log(JSON.stringify(response));
        callback(response);
    });

    stub.on('error', function(err){
        console.log(JSON.stringify(err));
        if(errHandler && typeof errHandler === "function") {
            errHandler(err);
        }
    });

    stub.on('end', function(){
        stub.end();
        if(endHandler && typeof endHandler === "function") {
            endHandler();
        }
    });

    requests.forEach(function(req){
        stub.write(req);
    });

    stub.end();
}

function newPriceRequest(trade) {
    const pricingdate = Math.round((new Date()).getTime() / 1000);  // always in UTC
    return {
      client_id: trade.contract.ticker,
      pricingdate: pricingdate,
      strike: trade.contract.strike,
      expiry: Math.round(Date.parse(trade.contract.expiry) / 1000),
      put_call: trade.contract.putcall,
      spot: trade.contract.index[0].quotes[0].close, // check if no quotes
      vol: 0.1,
      rate: 0.01
    };
  }
