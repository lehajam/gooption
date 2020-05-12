const koa = require('koa');
const fetch = require("node-fetch");
const grpc = require('grpc');
const protoLoader = require('@grpc/proto-loader');

// the app server
const app = new koa();
var pricer = newPricer();

setInterval(function(){
    console.log('Pricing all trades...');

    graphQuery(
      "http://localhost:8080/graphql",
      `{
        queryTrade {
          contract {
            ticker
            index {
              quotes (order: { desc: datePublished }, first: 1) {
                ... on StockQuote {
                  datePublished
                  close
                }
              }
            }
            ... on EuropeanContract {
              strike
              expiry
              putcall
            }
          }
        }
      }`,
      function(queryResponse) {

        const timestamp = (new Date()).toISOString();
        const pricingdate = Math.round((new Date()).getTime() / 1000);  // always in UTC
        const priceRequests = queryResponse.data.queryTrade.map(function (trade) {
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
        });

        grpcStreamQuery(
          pricer.Price(),
          priceRequests,
          function(response) {

            console.log(`contract:${response.client_id}\ttype:${response.value_type}\tvalue:${response.value}`);

            graphQuery(
              "http://localhost:8080/graphql",
              `mutation {
                addPriceResult(input: [{
                  datePublished:"${timestamp}",
                  contract: { ticker: "${response.client_id}" },
                  value: ${response.value},
                  resultType: "${response.value_type}",
                  source: "gobs"
                }]) {
                  priceresult {
                    id
                    value
                  }
              }}`);
          });
      }
    );
 }, 10000);


app.listen(4200);

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
  stub.on('data', callback);

  stub.on('error', function(err){
    console.json(JSON.stringify(err));
    if(errHandler && typeof errHandler === "function") {
      errHandler(err);
    }
  });

  stub.on('end', function(){
    if(endHandler && typeof endHandler === "function") {
      stub.end();
      endHandler();
    }
  });

  requests.forEach(function(req){
      stub.write(req);
  });

  stub.end();
}