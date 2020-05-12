
module.exports = {
    allTradesQuery,
    priceResultMutation
}

function allTradesQuery() {
    return `{
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
    }`;
}

function priceResultMutation(contractId, value, type, source) {
    const timestamp = (new Date()).toISOString();
    return `mutation {
        addPriceResult(input: [{
            datePublished:"${timestamp}",
            contract: { ticker: "${contractId}" },
            value: ${value},
            resultType: "${type}",
            source: "${source}"
        }]) {
            priceresult {
                id
                value
            }
        }
    }`;
}