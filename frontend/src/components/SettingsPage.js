import React, {useState, useEffect } from "react";

const configKeys = {
    dbPath: "",
    loConn: "",
    highConn: "",
    powDiff: ""
}

const SettingsPage = () => {
    const [settings, updateSettings] = useState(configKeys)
    const [loading, setLoading] = useState(true);
    const [status, setStatus] = useState("")

    const loadSettings = () => {
        window.backend.GetJsonConfigs().then( result => {
            var settingsMap = {};
            settingsMap["dbPath"] =   result.database["storage-path"]
            settingsMap["loConn"] =   result.network["connection-low"]
            settingsMap["highConn"] = result.network["connection-high"]
            settingsMap["powDiff"] =  (result.security["proofofwork-level"] / 4) - 4
            updateSettings(settingsMap)
            setLoading(false)
        })
    }

    const saveSettings = () => {
        window.backend.UpdateConfig(settings.dbPath, settings.loConn, 
            settings.highConn, settings.powDiff).then( res => {
                setStatus("successfully saved settings")
            })
    }

    useEffect(() => {
        loadSettings()
    }, [])

    return (
        <div>
        { !loading &&
        <div className="text-center">
            <h1>Settings</h1>
            <hr/>
            <div className="row justify-content-center mb-3">
                <div className="col-xl-8 mt-3">
                    <label htmlFor="dbPath" className="form-label">Local Storage Location</label>
                    <input type="text" className="form-control" id="dbPath" defaultValue={settings.dbPath} onChange={(e) => {
                        settings.dbPath = e.target.value
                    }}/>

                    <label htmlFor="loConn" className="form-label mt-3">Number of peers to connect to</label>
                    <div className="row">
                    <div className="col-2">Min: {settings.loConn}</div>
                    <div className="col-10"><input type="range" className="form-range" id="loConn" min="25" max="400" defaultValue={settings.loConn} onChange={(e) => {
                        updateSettings({...settings, loConn: parseInt(e.target.value)})
                    }}/></div>
                    </div>
                    
                    <div className="row">
                    <div className="col-2">Max: {settings.highConn}</div>
                    <div className="col-10"><input type="range" className="form-range" id="highConn" min="25" max="400" defaultValue={settings.highConn} onChange={(e) => {
                        updateSettings({...settings, highConn: parseInt(e.target.value)})
                    }}/></div>
                    </div>
                    {settings.loConn > settings.highConn && <p className="text-danger">minimum connections must not be higher than maximum connections</p>}

                    <label htmlFor="powDiff" className="form-label mt-3">Proof of Work Difficulty</label>
                    <input type="range" className="form-range" id="powDiff" min="0" max="3" defaultValue={settings.powDiff} onChange={(e) => {
                        settings.powDiff = parseInt(e.target.value)
                    }}/>

                    <button type="button" className="btn btn-primary" disabled={settings.loConn > settings.highConn} onClick={saveSettings}>Save Settings</button>
                    <h4 className="mt-3">{status}</h4>
                </div>
            </div>
        </div>
        }
        </div>
    );
}

export default SettingsPage