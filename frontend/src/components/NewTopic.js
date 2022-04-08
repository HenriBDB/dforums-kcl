import React, {useState } from "react";
import {getTitleValidationMessage, getDetailValidationMessage, 
    validateTitle, validateDetail} from "../util/Validation"; 

const NewTopic = () => {
    const [status, setStatus] = useState("")
    const [loading, setLoading] = useState(false)
    const [topic, setTopic] = useState("")
    const [detail, setDetail] = useState("")

    const createTopic = () => {
        setStatus("")
        setLoading(true)
        window.backend.ViewHandler.CreateTopic(topic, detail).then( res => {
                setLoading(false)
                setTopic("")
                setDetail("")
                setStatus("Successfully created topic")
            })
    }

    return (
        <div className="text-center">
        <h1>Create a New Topic</h1>
        <hr/>
        <div className="row justify-content-center mb-3">
            <div className="col-xl-8 mt-3">
                <label htmlFor="topic-field" className="form-label">
                    <h4>What is the title of your topic?</h4>
                </label>
                <div>
                <input type="text" className={"form-control " + (validateTitle(topic) ? "" : "is-invalid")} id="topic-field" value={topic} 
                aria-describedby="title-validation" onChange={(e) => setTopic(e.target.value)}/>
                <div id="title-validation" className="invalid-feedback">
                    {getTitleValidationMessage()}
                </div>
                </div>

                <div>
                <label htmlFor="detail-field" className="form-label mt-3">Can you add more detail about your topic? Do you have any references you would like to add?</label>
                <textarea className={"form-control " + (validateDetail(detail) ? "" : "is-invalid")} id="detail-field" rows="4" value={detail} 
                aria-describedby="detail-validation" onChange={(e) => setDetail(e.target.value)}/>
                <div id="detail-validation" className="invalid-feedback">
                    {getDetailValidationMessage()}
                </div>
                </div>
                
                <div className="mt-3">
                    <button className="btn btn-primary" 
                        disabled={loading || topic.length === 0 || detail.length === 0 || !validateTitle(topic) || !validateDetail(detail)} 
                        onClick={createTopic}>Create Topic</button>
                </div>
            </div>
        </div>
        { loading && 
            <div className="spinner-border" role="status">
                <span className="visually-hidden">Loading...</span>
            </div> }
        <h4>{status}</h4>
        </div>
    );
}

export default NewTopic;