import React from 'react';

class Content extends React.Component {
    render() {
        return (
            <div className="Content" style={{ padding: '16px' }}>
                {this.props.children}
            </div>
        );
    }
}

export default Content;
