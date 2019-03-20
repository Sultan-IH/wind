import React, { Component } from 'react';
import moment from 'moment';
import './App.css';

class App extends Component {
  constructor() {
    super()
    this.state = { lines: [] }
    this.appendLogs = this.appendLogs.bind(this)
    this.wss = []
    for (let i = 0; i < 10; i++) {
      let ws = new ReceiverWS(i, this.appendLogs);
      ws.start()
      this.wss.push(ws)
    }
  }

  appendLogs(line) {
    this.setState(state => {
      let lines = state.lines
      if (lines.length > 400) {
        lines.pop()
      }
      lines.unshift(line)
      return { lines }
    })
  }

  render() {
    return (
      <div className="App">
        <header className="App-header">
          <h3>Real time log data from the server</h3>
        </header>
        <div className="real-time-logs">
          {this.state.lines.map(line => {
            return <p> {"[" + moment().format() + "]" + line} </p>
          })
          }

        </div>
      </div>
    );
  }
}
class ReceiverWS {
  constructor(id, onmessage) {
    this.id = id
    this.start = this.start.bind(this)
    this.makeLine = this.makeLine.bind(this)
    this.onmessage = onmessage
  }

  start() {
    console.log("[WS:" + this.id + "] trying to connect ...")
    let ws = new WebSocket("ws://localhost:9009/receive/" + this.id);
    this.ws = ws
    this.ws.onopen = () => console.log("[WS:" + this.id + "] is opened")
    this.ws.onclose = () => {
      console.log("[WS:" + this.id + "] is closed; reconnecting ...")
      setTimeout(this.start, 5000)

    }
    this.ws.onmessage = (msg) => {
      console.log("[WS:" + this.id + "] got msg: ", msg)
      let line = this.makeLine(msg)
      this.onmessage(line)
    }

    this.ws.onmessage.bind(this)
  }
  makeLine(msg) {
    return "[WS:" + this.id + "]  wind sensor value is " + "[" + msg.data + "]"
  }

}
export default App;
