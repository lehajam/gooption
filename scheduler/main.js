const koa = require('koa');
const helpers = require('./helpers');
const queries = require('./queries');

var pricer = helpers.newPricer();
setInterval(function(){
    console.log('Pricing all trades...');
    // get all trades
    helpers.graphQuery("http://localhost:8080/graphql", queries.allTradesQuery(), function(queryResponse) {
        const priceRequests = queryResponse.data.queryTrade.map(helpers.newPriceRequest);
        helpers.grpcStreamQuery(pricer.Price(), priceRequests, function(response) {
            console.log(`contract:${response.client_id}\ttype:${response.value_type}\tvalue:${response.value}`);
            helpers.graphQuery("http://localhost:8080/graphql", queries.priceResultMutation(response.client_id, response.value, response.value_type, "gobs"));
          });
      }
    );
}, 10000);

const app = new koa();
app.listen(4200);
