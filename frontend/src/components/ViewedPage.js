import React from "react";
import { Link } from "react-router-dom";

const ViewedPage = ({ topics }) => {

    return (
        <div>
            <div className="text-center">
            <h1>Your Already Viewed Topics</h1>
            <hr/>
            {Object.keys(topics).length === 0 &&
                <div className="mt-3">
                    You haven't viewed any topics yet! Head over to the <Link to="/search">Search</Link> page to find topics for viewing.
                </div>
            }
            {<div className="row row-cols-1 row-cols-md-2 row-cols-xl-3">
                {Object.keys(topics).map( n => {
                    return(
                    <div key={n} className="col">
                        <div className="card text-center bg-light mb-3">
                            <div className="card-header"><strong>{topics[n].Short}</strong></div>
                            <div className="card-body">
                                <p className="card-text">{topics[n].Long}</p>
                                <Link to={"/search/"+n}>
                                    <button type="button" className="btn btn-primary mt-2">Explore</button>
                                </Link>
                            </div>
                        </div>
                    </div>);
                })}
            </div>}
            </div>
        </div>
    );
}

export default ViewedPage
