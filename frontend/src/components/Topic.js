import React, { useState } from 'react';
import { useParams } from "react-router-dom";
import {getTitleValidationMessage, getDetailValidationMessage, 
    validateTitle, validateDetail} from "../util/Validation"; 
import { indicatorToText } from '../util/Util';

function Topic({ data, loadMore, newComment }) {

    let params = useParams();

	return (
        <>
        { data && data[params.topicId] && 
        <><div className="text-center">
            <div className="alert alert-primary alert-dismissible fade show" role="alert">
                Note that comments are displayed in a random order
                <button type="button" className="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
            </div>
            <h1>{data[params.topicId].Short}</h1>
            <p>{data[params.topicId].Long}</p>
		</div>
        <hr/>
        <NewComment parent={params.topicId} newComment={newComment}/>
        <div>
        { data[params.topicId].Children && 
        <ul className="mt-3 p-0">
            { Object.keys(data[params.topicId].Children).map(k => {
                return (<Comment key={k} data={{...data[params.topicId].Children[k], ID: k}}
                    loadMore={loadMore} newComment={newComment}
                />);
            }) }
        </ul>
        }
        </div></>}
        </>
	);
}

export default Topic;

function Comment({ data = {}, loadMore, newComment }) {
    const [childrenLoaded, setChildrenLoaded] = useState(false)
    return (
        <div className="card text-left text-dark bg-light mt-3">
            <div className="card-header">
                {data.Short} - {indicatorToText(data.Indicator)}
            </div>
            <div className="card-body pb-2">
                <p className="card-text">{data.Long}</p>
                <hr/>
                { !childrenLoaded && <button className="btn btn-outline-dark me-2 py-1" type="button" onClick={() => {
                    loadMore(data.ID); setChildrenLoaded(true);
                }}>load more...</button>}
                <NewComment parent={data.ID} loadMore={loadMore} newComment={newComment}/>
                { data.Children && 
                <ul>
                    { Object.keys(data.Children).map(k => {
                        return (<Comment key={k} data={{...data.Children[k], ID: k}} loadMore={loadMore} newComment={newComment}/>);
                    }) }
                </ul>
                }
            </div>
        </div>
    );
}

function NewComment(props) {
    const [showModal, setShowModal] = useState(false);

    const [indicator, setIndicator] = useState(5)
    const [topic, setTopic] = useState("")
    const [detail, setDetail] = useState("")

    const cleanUpAndClose = () => {
        setShowModal(false)
        setIndicator(5)
        setTopic("")
        setDetail("")
    }

    return (
        <>
        { !showModal && <button className="btn btn-dark py-1" type="button" onClick={() => {
            setShowModal(true)
        }}>New Comment</button> }
        { showModal && 
        <form>
            <div className="col-xl-6">
                <div>
                <label htmlFor="topic-field" className="form-label">What is the topic of your comment?</label>
                <input type="text" className={"form-control " + (validateTitle(topic) ? "" : "is-invalid")} id="topic-field" value={topic} onChange={(e) => setTopic(e.target.value)}/>
                <div className="invalid-feedback">
                    {getTitleValidationMessage()}
                </div>
                </div>

                <label htmlFor="indicator-field" className="form-label mt-3">To what extent do you agree with the topic you chose? - {indicatorToText(indicator)}</label>
                <div className="row">
                    <div className="col-3 text-center">Disagree</div>
                    <div className="col-6">
                        <input type="range" className="form-range" id="indicator-field" min="0" max="10" defaultValue={indicator} onChange={(e) => {
                            setIndicator(e.target.value)
                        }}/>
                    </div>
                    <div className="col-3 text-center">Agree</div>
                </div>

                <div>
                <label htmlFor="detail-field" className="form-label mt-3">Why do you agree or disagree with what was said?</label>
                <textarea className={"form-control " + (validateDetail(detail) ? "" : "is-invalid")} id="detail-field" rows="4" value={detail} onChange={(e) => setDetail(e.target.value)}/>
                <div className="invalid-feedback">
                    {getDetailValidationMessage()}
                </div>
                </div>
                
                <div className="mt-3">
                    <button className="btn btn-outline-danger" onClick={() => {
                        cleanUpAndClose()
                    }}>Cancel</button>
                    <button className="btn btn-primary ms-2" 
                    disabled={topic.length === 0 || detail.length === 0 || !validateTitle(topic) || !validateDetail(detail)} 
                    onClick={() => {
                        props.newComment(topic, detail, indicator, props.parent)
                        cleanUpAndClose()
                    }}>Add Comment</button>
                </div>
            </div>
        </form>
        }
        </>
    )
}