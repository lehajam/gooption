import React from 'react';
import Link from '@material-ui/core/Link';
import {makeStyles} from '@material-ui/core/styles';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Title from './Title';
import gql from "graphql-tag";
import {useQuery} from "@apollo/react-hooks";
import moment from "moment";

// Generate Order Data
function createData(id, date, name, shipTo, paymentMethod, amount) {
  return {id, date, name, shipTo, paymentMethod, amount};
}

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

function preventDefault(event) {
  event.preventDefault();
}

const useStyles = makeStyles((theme) => ({
  seeMore: {
    marginTop: theme.spacing(3),
  },
}));

export default function Trades() {
  const classes = useStyles();
  const {data, loading, error} = useQuery(GET_TRADES);
  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error</p>;

  return (
      <React.Fragment>
        <Title>Trades</Title>
        <Table size="small">
          <TableHead>
            <TableRow>
              <TableCell>Trade Date</TableCell>
              <TableCell>Quantity</TableCell>
              <TableCell>Contract</TableCell>
              <TableCell>Underlying</TableCell>
              <TableCell>Spot</TableCell>
              <TableCell>Strike</TableCell>
              <TableCell>Moneyness</TableCell>
              <TableCell>Expiry</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {data.queryTrade.map((row) => (
                <TableRow key={row.id}>
                  <TableCell type="date">{moment(row.tradeDate).format('D/MM/YYYY')}</TableCell>
                  <TableCell>{row.quantity}</TableCell>
                  <TableCell>{row.contract.symbol}</TableCell>
                  <TableCell>{row.contract.refIndex[0].symbol}</TableCell>
                  <TableCell>{row.contract.refIndex[0].quotes && row.contract.refIndex[0].quotes[0].last}</TableCell>
                  <TableCell>{row.contract.strike}</TableCell>
                  <TableCell>{row.contract.refIndex[0].quotes && row.contract.strike / row.contract.refIndex[0].quotes[0].last}</TableCell>
                  <TableCell>{moment(row.contract.expiry).format('D/MM/YYYY')}</TableCell>
                </TableRow>
            ))}
          </TableBody>
        </Table>
        <div className={classes.seeMore}>
          <Link color="primary" href="#" onClick={preventDefault}>
            See more orders
          </Link>
        </div>
      </React.Fragment>
  );
}
