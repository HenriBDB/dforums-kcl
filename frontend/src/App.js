import React from 'react';
import { MemoryRouter, Routes, Route } from "react-router-dom";
import './App.css';
import MenuBar from './components/MenuBar';
import Topic from './components/Topic';
import SettingsPage from './components/SettingsPage';
import SearchPage from './components/SearchPage';
import NewTopic from './components/NewTopic';
import ViewedPage from './components/ViewedPage';
import { addTopic, addNode, addNodes } from './util/DataHandler';

class App extends React.Component {
  state = {
    topics: [],
    data: {},
    inProgress: [],
  }

  newComment = (topic, detail, indicator, parent) => {
    this.setState({...this.state, inProgress: [...this.state.inProgress, "comment"]})
    window.backend.ViewHandler.CreateNode(topic, detail, parseInt(indicator), parent).then(() => {
      this.state.inProgress.shift()
      this.setState({...this.state, inProgress: [...this.state.inProgress]})
    })
  }

  componentDidMount() {
    window.wails.Events.On('new_node', node => {
      this.receiveNode(node)
    })
  }

  getParentNodes = () => {
    window.backend.ViewHandler.GetAllTopics().then( res => {
      this.setState({...this.state, topics: res})
    })
  }

  receiveNode = (node) => {
    const newState = Object.assign({}, this.state.data);
    addNode(node, newState)
    this.setState({...this.state, data: newState})
  }

  registerChildren = (nodeID) => {
    window.backend.ViewHandler.GetChildren(nodeID).then( children => {
      const newState = Object.assign({}, this.state.data);
      addNodes(children, newState)
      this.setState({...this.state, data: newState})
    })
  }

  registerTopic = (topic) => {
    if (!this.state.data[topic.ID]) {
      window.backend.ViewHandler.GetChildren(topic.ID).then( children => {
        const newState = Object.assign({}, this.state.data);
        addTopic(topic, newState)
        addNodes(children, newState)
        this.setState({...this.state, data: newState})
      })
    }
  }

  render() {
    return <MemoryRouter>
      <div id="app" className="App bg-light">
        <MenuBar/>
        <div className="main-frame">
          <div className="p-4 mb-5 justify-content-center">
            <div className="container">
              <Routes>
                <Route path="search" element={
                  <SearchPage getTopics={this.getParentNodes} register={this.registerTopic} topics={this.state.topics}/>
                } exact />
                <Route path="search/:topicId" element={<Topic data={this.state.data} loadMore={this.registerChildren} newComment={this.newComment}/>} exact />
                <Route path="new-topic" element={
                  <NewTopic/>
                }/>
                <Route path="settings" element={
                  <SettingsPage/>
                }/>
                <Route path="viewed" element={
                  <ViewedPage topics={this.state.data}/>
                }/>
                <Route path="*" element={
                  <SearchPage getTopics={this.getParentNodes} register={this.registerTopic} topics={this.state.topics}/>
                }/>
              </Routes>
            </div>
          </div>
          <footer className="text-center text-white w-100">
            <div className="p-3">
                {this.state.inProgress.map((e, i) => {
                  return <span key={i} className="badge rounded-pill bg-success d-inline-flex flex-row align-items-center mx-2">
                    <div>creating {e}&nbsp;&nbsp;</div><div className="spinner-grow spinner-grow-sm text-light" role="status">
                    <span className="visually-hidden">Loading...</span>
                  </div></span>
                })}
            </div>
          </footer>
        </div>
      </div>
    </MemoryRouter>;
  }
}

export default App;
