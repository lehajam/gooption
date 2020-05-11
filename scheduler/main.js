const koa = require('koa');
const fetch = require("node-fetch");
const grpc = require('grpc');
const protoLoader = require('@grpc/proto-loader');

// the app server
const app = new koa();
var pricer = newPricer();

setInterval(function(){
    console.log('Pricing all trades...');

    fetch("https://c14f96ce.ngrok.io/graphql", {
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
          console.log(
            "contract:" + response.client_id + '\t' +
            "type: " + response.value_type + '\t' +
            "value: " + response.value);
        });

        price.on('error', function(e) {
            console.log(JSON.stringify(e));
        });

        price.on('end', function() {
            console.log('Trades have been priced successfully');
            price.end();
        });

        const now = new Date();
        const pricingdate = Math.round(now.getTime() / 1000);  // always in UTC
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