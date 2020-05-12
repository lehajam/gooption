const koa = require('koa');
const fetch = require("node-fetch");
const grpc = require('grpc');
const protoLoader = require('@grpc/proto-loader');

// the app server
const app = new koa();
var pricer = newPricer();

setInterval(function(){
    console.log('Pricing all trades...');

    fetch("http://localhost:8080/graphql", {
        "method": "POST",
        "headers": {
          "content-type": "application/json"
        },
        "body": JSON.stringify({ query: `{
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
          }` })
    })
    .then(res => res.json())
    .then(jsonres => {

        var price = pricer.Price();

        price.on('data', function(response) {

          var timestamp = (new Date()).toISOString();
          console.log(`contract:${response.client_id}\ttype:${response.value_type}\tvalue:${response.value}`);
          fetch("http://localhost:8080/graphql", {
            "method": "POST",
            "headers": {
              "content-type": "application/json"
            },
            "body": JSON.stringify({ query: `mutation {
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
            }}`})
          })
          .then(res => res.json())
          .then(json => {
            console.log('inserted row');
          });
        });

        price.on('error', function(e) {
            console.log(JSON.stringify(e));
        });

        price.on('end', function() {
            console.log('Trades have been priced successfully');
            price.end();
        });

        const pricingdate = Math.round((new Date()).getTime() / 1000);  // always in UTC
        jsonres.data.queryTrade.forEach(function(trade){
            price.write({
                client_id: trade.contract.ticker,
                pricingdate: pricingdate,
                strike: trade.contract.strike,
                expiry: Math.round(Date.parse(trade.contract.expiry) / 1000),
                put_call: trade.contract.putcall,
                spot: trade.contract.index[0].quotes[0].close, // check if no quotes
                vol: 0.1,
                rate: 0.01
            });
        });

        price.end();
    })
    .catch(err => {
        console.log(err);
    });
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