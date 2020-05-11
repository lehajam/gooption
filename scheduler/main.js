const koa = require('koa');
const phin = require('phin');
const grpc = require('grpc');
const protoLoader = require('@grpc/proto-loader');


// the app server
const app = new koa();
var pricer = newPricer();

app.use(async ctx => {

    var price = pricer.Price();

    price.on('data', function(response) {
      console.log("type: " + response.value_type + '\t' + "value: " + response.value);
    });

    price.on('error', function(e) {
        console.log(JSON.stringify(e));
    });

    price.on('end', function() {
        price.end();
    });

    price.write({
        pricingdate: 1564669224,
        strike: 159.76,
        expiry: 1582890485,
        put_call: 'put',
        spot: 264.16,
        vol: 0.1,
        rate: 0.01
    });

    price.end();
});

console.log('server ready on localhost:4200');
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