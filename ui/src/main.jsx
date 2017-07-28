import React from 'react'
import { render } from 'react-dom'
import { Router, Route, IndexRoute, Link } from 'react-router';
import injectTapEventPlugin from 'react-tap-event-plugin';
import App from './components/App';
import Dashboard from './components/Dashboard';
import Promise from 'promise-polyfill';
import 'whatwg-fetch';

// To add to window
if (!window.Promise) {
    window.Promise = Promise;
}

injectTapEventPlugin();



render((
    <Router>
        <Route path="/" component={App}>
            <IndexRoute component={Dashboard} />
        </Route>
    </Router>
), document.getElementById('app'));
