import React from 'react'
import getMuiTheme from 'material-ui/styles/getMuiTheme';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
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


import Content from './Content';
import Footer from './Footer';
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

class Dashboard extends React.Component {
    constructor(props, context) {
        super(props, context);

        this.state = {
            status: null,
        }

        this.openHistory = this.openHistory.bind(this);
        this.openDashboard = this.openDashboard.bind(this);
        this.loadStatus = this.loadStatus.bind(this);
        this.lastEtag = null;
    }

    componentDidMount() {
        this.loadStatus();
        this.interval = setInterval(this.loadStatus, this.props.pollInterval);
    }

    componentWillUnmount() {
        clearInterval(this.interval);
    }

    loadStatus() {
        var that = this;
        var headers = new Headers();
        if (this.lastEtag !== null) {
            headers.append("If-None-Match", this.lastEtag);
        }
        fetch(`${BASE_PATH}/api/status`, { credentials: 'same-origin', headers: headers })
            .then(function (response) {
                var etag = response.headers.get('ETag');
                if (etag !== that.lastEtag) {
                    that.lastEtag = etag;
                    return response.json();
                }
                return null;
            }).then(function (json) {
                if (json === null) {
                    return;
                }
                that.setState({ status: json });
            }).catch(function (ex) {
                console.log('failed to fetch status', ex)
            });
    }

    openHistory() {
        history.push(`/history`);
    }

    openDashboard() {
        history.push('/');
    }

    render() {
        return (
            <div className="dashboard">
                <AppBar title=""
                    title={<span style={{ cursor: 'pointer' }}>Service Manager</span>}
                    onTitleTouchTap={this.openDashboard}
                    showMenuIconButton={false}
                    iconElementRight={
                        <IconMenu iconButtonElement={<IconButton><MoreVertIcon /></IconButton>}
                            targetOrigin={{ horizontal: 'right', vertical: 'top' }}
                            anchorOrigin={{ horizontal: 'right', vertical: 'top' }}>
                            <MenuItem primaryText="History" onClick={this.openHistory} />
                            <MenuItem primaryText="Log out" />
                        </IconMenu>} />
                <Content>
                    <Table selectable={false}>
                        <TableHeader displaySelectAll={false} adjustForCheckbox={false}>
                            <TableRow>
                                <TableHeaderColumn>Service</TableHeaderColumn>
                                <TableHeaderColumn>Status</TableHeaderColumn>
                                <TableHeaderColumn>Details</TableHeaderColumn>
                                <TableHeaderColumn></TableHeaderColumn>
                            </TableRow>
                        </TableHeader>
                        <TableBody displayRowCheckbox={false}>
                            {this.renderServices()}
                        </TableBody>
                    </Table>
                    <Footer />
                </Content>
            </div>
        );
    }

    renderServices() {
        if (this.state.status === null) {
            return null;
        }
        var status = this.state.status;
        var result = Array();
        for (let provider of Object.keys(status)) {
            if (status[provider].success === false) {
                continue;
            }
            status[provider].services.forEach(function (s) {
                result.push(<Service
                    key={`${provider}__${s.service}`}
                    provider={provider}
                    service={s.service}
                    url={s.url}
                    status={s.status}
                    detail={s.detail}
                    can_start={s.can_start}
                    can_stop={s.can_stop}
                    can_restart={s.can_restart} />);
            });
        }
        return result;
    }
}

Dashboard.propTypes = {
    pollInterval: React.PropTypes.number.isRequired
};

Dashboard.defaultProps = {
    pollInterval: 1000
};

export default Dashboard;
