import React from 'react';
import './App.css';
import {useQuery} from '@apollo/react-hooks';
import gql from "graphql-tag";

const GET_TRADES = gql`
    {
        queryTrade {
            id
            tradeDate
            quantity
            contract {
                symbol
                refIndex {
                    symbol
                    quotes(order: { desc: datePublished}, first: 1) {
                        symbol
                        datePublished
                        last
                    }
                }
                ... on Option {
                    strike
                    expiry
                    putcall
                    optionType
                }
            }
            valuations(order: { desc: datePublished}, first: 1) {
                ... on Price {
                    datePublished
                    value
                }
            }
        }
    }`;


function App() {
    const {data, loading, error} = useQuery(GET_TRADES);
    if (loading) return <p>Loading...</p>;
    if (error) return <p>Error</p>;

    return (
        <React.Fragment>
            <h1>Trades</h1>
            <div className="container">
                {data &&
                data.queryTrade &&
                data.queryTrade.map((trade, index) => (
                    <div key={index} className="card">
                        <div class="card-body">
                            <h3>{trade.id}</h3>
                            <h3>{trade.contract.strike}</h3>
                            <h3>{trade.contract.expiry}</h3>
                            <p>
                                {trade.valuations && trade.valuations.length !== 0 && (
                                    <p>
                                        {" "}
                                        Prices:
                                        {trade.valuations.map((v, indx) => {
                                            return <p key={indx}> {v.value} </p>;
                                        })}
                                    </p>
                                )}
                            </p>
                        </div>
                    </div>
                ))}
            </div>
        </React.Fragment>
    );

}

export default App;
