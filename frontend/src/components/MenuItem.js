import React from 'react';
import { Link } from "react-router-dom";

function MenuItem(props) {

    return(
        <Link to={props.link}>
            <button type="button" className="btn btn-secondary w-100 rounded-0">
                {props.menuName}
            </button>
        </Link>
    );

}

export default MenuItem;
