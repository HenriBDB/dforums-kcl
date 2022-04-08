import _ from 'lodash';

const pathMap = new Map()

export const addTopic = (topic, state) => {
    pathMap.set(topic.ID, "")
    state[topic.ID] = { Short: topic.Short, Long: topic.Long, 
        Indicator: topic.Indicator, Children: {} }
}

export const addNodes = (nodes, state) => {
    nodes.forEach(n => {
        addNode(n, state)
    })
}

export const addNode = (node, state) => {
    // assert no topics
    if (pathMap.has(node.Parent)) {
        // Add node to path map
        var parentPath = [...pathMap.get(node.Parent)]
        parentPath.push(node.Parent)
        pathMap.set(node.ID, parentPath)
        // Update state
        var lodashPath = [];
        for (var i=0; i<pathMap.get(node.ID).length; ++i) {
            lodashPath.push(pathMap.get(node.ID)[i])
            lodashPath.push("Children")
        }
        lodashPath.push(node.ID)
        _.set(state, 
            lodashPath, 
            {Short: node.Short, Long: node.Long, Indicator: node.Indicator, Children: {}}
        )
    }
}