import React from 'react'
import getMuiTheme from 'material-ui/styles/getMuiTheme';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import { Link } from 'react-router';
import AppBar from 'material-ui/AppBar';
import { Table, TableHeader, TableHeaderColumn, TableRow, TableRowColumn, TableBody } from 'material-ui/Table';
import FlatButton from 'material-ui/FlatButton';
import IconButton from 'material-ui/IconButton';
import FontIcon from 'material-ui/FontIcon';
import IconMenu from 'material-ui/IconMenu';
import MenuItem from 'material-ui/MenuItem';
import Dialog from 'material-ui/Dialog';
import AVStop from 'material-ui/svg-icons/av/stop';
import MoreVertIcon from 'material-ui/svg-icons/navigation/more-vert';


import Service from './Service';

const muiTheme = getMuiTheme({
    tooltip: {
        zIndex: 9999,
    },
});

const appLinkStyle = {
    color: 'white',
    textDecoration: 'none',
    cursor: 'pointer',
};

class App extends React.Component {

    render() {
        return (
            <MuiThemeProvider muiTheme={muiTheme}>
                {this.props.children}
            </MuiThemeProvider>
        );
    }

}

export default App;
