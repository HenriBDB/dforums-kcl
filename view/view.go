package view

import (
	"dforum-app/configuration"
	"dforum-app/security"
	"dforum-app/storage"
	"encoding/base64"
	"math/rand"

	bloom "github.com/bits-and-blooms/bloom/v3"
	"github.com/wailsapp/wails"
)

/*
	This file acts as a link between the presentation and the data layers
*/

type GuiNode struct {
	ID        string
	Parent    string
	Short     string
	Long      string
	Indicator int
}

type ViewHandler struct {
	storageModule *storage.StorageModule
	// Remember all hashes registered on the GUI -> uses the base hash and NOT base64URL ones
	filter       *bloom.BloomFilter
	wailsRuntime *wails.Runtime
}

func (vh *ViewHandler) WailsInit(runtime *wails.Runtime) error {
	vh.wailsRuntime = runtime
	return nil
}

func NewViewHandler(storage *storage.StorageModule) *ViewHandler {
	vh := &ViewHandler{
		storageModule: storage,
		filter:        bloom.NewWithEstimates(10000, 0.01),
	}
	// Subscribe to storage to be updated with incoming nodes from the network
	// Enables real time updating of the GUI
	storage.Subscribe(vh)
	return vh
}

func (vh *ViewHandler) GetAllTopics() []GuiNode {
	parentNodes := vh.storageModule.GetTopLevelNodes()
	return nodesToGuiNodes(parentNodes)
}

func (vh *ViewHandler) CreateTopic(topic string, detail string) {
	newTopic := storage.NewNode(topic, detail, -1, security.HashSignature{})
	vh.storageModule.StoreAndRegisterNewNode(newTopic)
}

func (vh *ViewHandler) CreateNode(topic string, detail string, indicator int, parent string) {
	newNode := storage.NewNode(topic, detail, int8(indicator), hashFromBase64(parent))
	vh.storageModule.StoreAndRegisterNewNode(newNode)
}

func (vh *ViewHandler) GetChildren(base64Id string) []GuiNode {
	hashId := hashFromBase64(base64Id)
	// Register parent when children are fetched
	vh.filter.Add(hashId[:])
	childrenNodes := vh.storageModule.GetChildrenNodes(hashId, false, -1)
	return nodesToGuiNodes(childrenNodes)
}

func (vh *ViewHandler) RegisterNewNode(node *storage.Node) {
	if node == nil {
		return
	}
	hashId := node.DatObj.Parent
	if vh.filter.Test(hashId[:]) {
		guiNode := convertNode(node)
		vh.wailsRuntime.Events.Emit("new_node", guiNode)
	}
}

// Converts storage nodes into a GUI compatible data struct
// Shuffles the resulting array to avoid always showing nodes in the same order
func nodesToGuiNodes(nodes []*storage.Node) []GuiNode {
	guiNodes := []GuiNode{}
	for _, v := range nodes {
		if v == nil {
			continue
		}
		guiNodes = append(guiNodes, convertNode(v))
	}
	// https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle
	for i := range guiNodes { // Shuffle
		j := rand.Intn(i + 1)
		guiNodes[i], guiNodes[j] = guiNodes[j], guiNodes[i]
	}
	return guiNodes
}

func convertNode(node *storage.Node) GuiNode {
	if node == nil {
		return GuiNode{}
	}
	return GuiNode{
		ID:        base64.URLEncoding.EncodeToString(node.SecObj.Fingerprint[:]),
		Parent:    base64.URLEncoding.EncodeToString(node.DatObj.Parent[:]),
		Short:     node.DatObj.Topic,
		Long:      node.DatObj.Content,
		Indicator: int(node.DatObj.Indicator),
	}
}

func hashFromBase64(base64Id string) security.HashSignature {
	if base64Id == "" { // For empty parent hashes signifying top level node
		return security.HashSignature{}
	}
	hash, err := base64.URLEncoding.DecodeString(base64Id)
	if err != nil {
		configuration.Logger.Error("failed to convert GUI ID to local ID")
		return security.HashSignature{}
	}
	return *(*[28]byte)(hash)
}
