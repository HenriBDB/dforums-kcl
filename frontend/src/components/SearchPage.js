import React, {useState, useEffect } from "react";
import { Link } from "react-router-dom";

const SearchPage = ({ getTopics, register, topics }) => {
    const [loaded, setLoaded] = useState(false)

    useEffect(() => {
        getTopics()
        setLoaded(true)
    }, [getTopics])

    return (
        <div>
            <div className="text-center">
            <h1>Find Topics</h1>
            <hr/>
            </div>
            <div className="alert alert-primary alert-dismissible fade show" role="alert">
                Note that topics are displayed in a random order
                <button type="button" className="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
            </div>
            { loaded && 
                <div className="row row-cols-1 row-cols-md-2 row-cols-xl-3">
                    {topics.map( n => {
                        return(
                        <div key={n.ID} className="col">
                            <div className="card text-center bg-light mb-3">
                                <div className="card-header"><strong>{n.Short}</strong></div>
                                <div className="card-body">
                                    <p className="card-text">{n.Long}</p>
                                    <Link to={"/search/"+n.ID}>
                                        <button type="button" className="btn btn-primary mt-2" onClick={() => {
                                            register(n)
                                        }}>Explore</button>
                                    </Link>
                                </div>
                            </div>
                        </div>);
                    })}
                </div>
            }
            <button className="btn btn-primary" onClick={getTopics}>Refresh</button>
        </div>
    );
}

export default SearchPage