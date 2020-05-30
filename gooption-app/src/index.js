import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import * as serviceWorker from './serviceWorker';

import {ApolloClient} from 'apollo-client';
import {InMemoryCache} from 'apollo-cache-inmemory';
import {HttpLink} from 'apollo-link-http';
import {ApolloProvider} from '@apollo/react-hooks';
import Dashboard from "./Dashboard";

const cache = new InMemoryCache();
const link = new HttpLink({
  uri: 'http://localhost:8080/graphql'
});

const client = new ApolloClient({
  cache,
  link
});

ReactDOM.render(<ApolloProvider client={client}><Dashboard/></ApolloProvider>, document.getElementById('root'));

serviceWorker.unregister();
