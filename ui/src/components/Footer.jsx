import React from 'react'
import { grey300 } from 'material-ui/styles/colors';

class Footer extends React.Component {

    render() {
        return (
            <div className="footer" style={{ marginTop: '32px', fontSize: '.9em', color: grey300 }}>
                © kervehrt 2017
      </div>
        );
    }

}

export default Footer;
