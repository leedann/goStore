import React, { Component } from 'react';
import './App.css';

const BASE_URL = "http://138.68.249.59/";
const SEARCHURL = "v1/summary?url=";

class OgpCard extends Component {
    render() {
        return (<div className="card">
            {this.props.image ? <img src={this.props.image} alt="ogp" /> : ""}
            <div className="card-content">
                {this.props.title ? <p> {this.props.title} </p> : ""}
                {this.props.url ? <p> {this.props.url} </p> : ""}
                {this.props.description ? <p> {this.props.description} </p> : ""}
            </div>
        </div>);
    }
}
class OgpSection extends Component {
  render() {
    var ogpAttr = this.props.ogp 
    if (!this.props.errorMess) {
        var cards = <OgpCard title={ogpAttr.title} image={ogpAttr.image} description={ogpAttr.description}>{ogpAttr.title}</OgpCard>
        return <div>{cards}</div>
    }else {
        return <p>{this.props.errorMess}</p>;
    }
  }
}


class App extends Component {
    constructor(props) {
        super(props);
        this.state = {
            query: "", 
            ogp: [],
            errorMess: ""
        };
    }

    handleSubmit(event) {
        event.preventDefault();
        fetch(BASE_URL + SEARCHURL + this.state.query)
            .then(response => response.json())
            .then(data => this.setState({
                ogp: data,
                errorMess: ""}))
            .catch(error =>  this.setState({errorMess: "Bad request from URL (try http or https)"})); 
    }
    handleChange(event) {
        this.setState({query: event.target.value});

    }

    render() {
        return (
            <div className="content-wrapper">
                <div className="App-header">
                  <h2>Welcome to OGP finder</h2>
                </div>
                <div className="ogpContent">
                  <form className="search-form"
                      onSubmit={event => this.handleSubmit(event)}>
                      <div className="input-group">
                          <input type="text" className="form-control"
                              value={this.state.query} 
                              placeholder="enter a valid url"
                              autoFocus                              
                              required
                              onChange={event => this.handleChange(event)} />
                          <span className="input-group-btn">
                              <button className="btn btn-primary" 
                                  aria-label="url search">
                                  Search
                              </button>
                          </span>
                      </div>
                  </form>
                </div>
                <div id="tile-wrapper">
                    <OgpSection ogp={this.state.ogp} errorMess={this.state.errorMess} />
                </div>
            </div>
        );
    }
}

export default App;
